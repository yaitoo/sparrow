package db

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/yaitoo/sparrow/db/parser"
)

type Context struct {
	sync.Mutex
	currendDbLock   sync.Mutex
	ctx             context.Context
	currentDb       *sql.DB
	currentDbID     int64
	currentConnStr  string
	tx              *sql.Tx
	database        *Database
	txStart         bool
	slowSQLDuration time.Duration
	// slowSQLDuration *time.Duration
}

//NewQuery create a new
func (c *Context) NewQuery(ctx context.Context) *Querys {

	return &Querys{
		// Context:   ctx,
		vars:      make(map[string]interface{}),
		hasRawSql: false,
		ctx:       c,
		// conditions: make(map[string][]string),
	}
}

//NewCommand create a new
func (c *Context) NewCommand(ctx context.Context) *Commands {
	return &Commands{
		Querys: Querys{
			hasRawSql: false,
			vars:      make(map[string]interface{}),
			ctx:       c,
		},
		kvPair: make(map[string]interface{}),
	}
}

//TryOpenDb 尝试打开db连接， 建议sql.DB的示例，如果已经打开， 则使用已打开的sql.DB实例，每个库一个sql.DB
func (c *Context) TryOpenDb(conn string, maxLifetime time.Duration, maxConns, minConns int) (*sql.DB, error) {
	db, err := c.database.connect(conn, maxLifetime, maxConns, minConns)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//Begin 启动事务，但是因为需要SQL来判断是要在哪一个库上启动事务，所以这里先做标志，等到第一个Command执行的时候，再执行tx的初始化
func (c *Context) Begin() error {
	c.Lock()
	defer c.Unlock()
	//defer cmd.ctx.Unlock()
	if c.tx != nil || c.txStart == true {
		return ErrInvalidTransaction
	}
	c.txStart = true
	return nil
}

func (c *Context) Commit() error {
	c.Lock()
	defer func() {
		c.tx = nil
		c.currentDb = nil
		c.currentDbID = -1
		c.txStart = false
		c.Unlock()
	}()
	if c.tx == nil {
		return ErrInvalidTransaction
	}

	if err := c.tx.Commit(); err != nil {
		// q.tx.Rollback()

		return ErrInvalidTransaction
	}

	return nil
}

func (c *Context) Rollback() error {
	c.Lock()
	defer func() {
		c.tx = nil
		c.currentDb = nil
		c.currentDbID = -1
		c.txStart = false
		c.Unlock()
	}()
	if c.tx == nil {
		return nil
	}

	if err := c.tx.Rollback(); err != nil {
		return ErrInvalidTransaction
	}

	return nil
}

// func (c *Context) SetCurrentDB(db *sql.DB, dbID int64) error {

// 	//c.Lock()
// 	//defer c.Unlock()

// 	if c.currentDb != nil && c.currentDbID != dbID {
// 		return errCrossDatabase
// 	}

// 	c.currentDb = db
// 	c.currentDbID = dbID

// 	tx, err := c.currentDb.BeginTx(c.ctx, &sql.TxOptions{
// 		Isolation: sql.LevelSerializable,
// 	})

// 	if err != nil {
// 		logger.Error(err)
// 		//return ErrInvalidTransaction
// 		return err
// 	}

// 	if c.tx != nil {
// 		return nil
// 	}

// 	c.tx = &Tx{tx: tx}
// 	return nil
// }

//PrepareContext 生成mysql的Statment, 并返回解析后的SQL和变数列表
func (c *Context) PrepareContext(ctx context.Context, commndText string, vars map[string]interface{}, tableName string) (*sql.Stmt, string, []interface{}, error) {
	var (
		cmd     parser.SqlStringAndEnitityData
		connStr string
		err     error
	)
	// get conn by config
	config := c.database.config()
	if cmd, err = parser.ParseWithVarMap(config, commndText, vars); err != nil {
		return nil, "", nil, err
	}
	entity := cmd.GetTarget()
	mdb := config.GetDatabase(entity.Key.Value.AlgVer, entity.Key.Value.DbId, entity.EntityName)
	if connStr, err = mdb.ConnStr(); err != nil {
		logger.Error(err)
		return nil, "", nil, err
	}
	var db *sql.DB
	if db, err = c.TryOpenDb(connStr, mdb.MaxLifeTIme, mdb.MaxConns, mdb.MinConns); err != nil {
		logger.Error(err)
		return nil, "", nil, err
	}

	//已经启动事务，判断当前SQL分库分表后等到的库位置，如果和上一个事务操作不在同一个库，则抛出异常，不允许执行。
	//因为已经在两个不同的库上，无法保证一致性事务。
	if c.tx != nil && c.currentDbID != entity.Key.Value.DbId {
		return nil, "", nil, errCrossDatabase
	}
	//事务已经被标记启动
	if c.txStart == true && c.tx == nil {

		c.currentDb = db
		c.currentDbID = entity.Key.Value.DbId

		tx, err := c.currentDb.BeginTx(c.ctx, &sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

		if err != nil {
			logger.Error(err)
			//return ErrInvalidTransaction
			return nil, "", nil, err
		}

		c.tx = tx

	}

	var stmt *sql.Stmt
	//事务已经启动，则统一使用事务的连接
	if c.tx != nil {
		stmt, err = c.tx.PrepareContext(ctx, cmd.SqlString)
	} else { //事务未启动， 使用瞬间的db连接
		stmt, err = db.PrepareContext(ctx, cmd.SqlString)
	}

	if err != nil {
		logger.Error(err)
		//return ErrInvalidTransaction
		return nil, "", nil, err
	}

	return stmt, cmd.SqlString, cmd.PrepareArgs, nil

}
