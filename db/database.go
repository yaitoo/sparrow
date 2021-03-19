package db

//http://go-database-sql.org/overview.html

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/yaitoo/sparrow/db/dbconfig"
	"github.com/yaitoo/sparrow/db/model"
	"github.com/yaitoo/sparrow/db/shardingId"
	"github.com/yaitoo/sparrow/log"
)

const fileLocation dbconfig.FileLocation = "./conf.d/db.yaml"

var (
	//InstanceID ID生成器实例ID，用于分散式生成时候保持ID唯一，可以使用以下代码做环境设定
	// func init(){
	//	db.InstanceID = types.Atoi(conf.Value("server","id","1"),1)
	//}
	InstanceID int64
	//Shardings 表ID生器成集合，用于保存各个表ID生成器的当前状态
	Shardings = make(map[string]map[string]*shardingId.IdGenerator)

	logger = log.NewLogger("db")
)

func init() {
	dbconfig.Initonfig(fileLocation)
}

//Database 数据库操作实例，包含ID生成状态，共享的db连线实例，设定文档
type Database struct {
	items     map[string]*sql.DB
	shardings map[string]map[string]*shardingId.IdGenerator
	config    func() model.Config //model.Config
	sync.Mutex
}

//NewDatabase  创建数据库操作实例
//d := db.NewContext(ctx,  db.WithConfig(config))
func NewDatabase(ctx context.Context, opts ...Option) *Database {

	d := &Database{
		items:     make(map[string]*sql.DB),
		shardings: make(map[string]map[string]*shardingId.IdGenerator),
		//config:    func() model.Config { return dbconfig.GetConfigObj() },
	}

	//执行自定义属性设定
	for _, opt := range opts {
		opt(d)
	}

	if d.config == nil {
		d.config = func() model.Config { return dbconfig.GetConfigObj() }
	}

	return d
}

//Open 创建一个数据库操作上下文对象
func (d *Database) Open(ctx context.Context) *Context {
	return &Context{
		database: d,
		ctx:      ctx,
	}
}

//Open create a database connection
func (d *Database) connect(conn string, maxLifetime time.Duration, maxConns, minConns int) (*sql.DB, error) {
	// create critical section
	// defer unlock
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	_db, ok := d.items[conn]
	if ok == false {
		_db, err := sql.Open("mysql", conn)
		if err != nil {
			return nil, errOpenConnError(conn, err)
		}
		// set conn pool
		_db.SetMaxOpenConns(maxConns)
		_db.SetMaxIdleConns(minConns)
		_db.SetConnMaxLifetime(maxLifetime)
		err = _db.Ping()
		if err != nil {
			return nil, errOpenConnError(conn, err)
		}
		d.items[conn] = _db
		return _db, nil
	}

	return _db, nil

	//_db, err := sql.Open()
	// step 1: open db
	// step 2 : add to map
	// step 3 : return db
}

func (d *Database) NewID(ctx context.Context, tableName, tag string) (int64, error) {
	config := d.config()
	d.Lock()
	defer d.Unlock()

	if tableShardings, ok := d.shardings[tableName]; ok {
		if tagSharding, ok := tableShardings[tag]; ok {
			return tagSharding.NextID()
		} else {
			generator, err := shardingId.NewIdGenerator(tableName, tag, config.GetNewestVersion(), InstanceID)
			if err != nil {
				return 0, err
			}
			d.shardings[tableName][tag] = generator
			return generator.NextID()
		}
	} else {
		d.shardings[tableName] = make(map[string]*shardingId.IdGenerator)
		generator, err := shardingId.NewIdGenerator(tableName, tag, config.GetNewestVersion(), InstanceID)
		if err != nil {
			return 0, err
		}
		d.shardings[tableName][tag] = generator
		return generator.NextID()
	}
}

func (d *Database) NewSubID(ctx context.Context, tableName, tag string, parendID int64) (int64, error) {
	config := d.config()
	d.Lock()
	defer d.Unlock()
	if tableShardings, ok := d.shardings[tableName]; ok {
		if tagSharding, ok := tableShardings[tag]; ok {
			return tagSharding.NextID()
		} else {
			generator, err := shardingId.NewSubId(tableName, tag, config.GetNewestVersion(), parendID)
			if err != nil {
				return 0, err
			}
			d.shardings[tableName][tag] = generator
			return generator.NextID()
		}
	} else {
		d.shardings[tableName] = make(map[string]*shardingId.IdGenerator)
		generator, err := shardingId.NewSubId(tableName, tag, config.GetNewestVersion(), parendID)
		if err != nil {
			return 0, err
		}
		d.shardings[tableName][tag] = generator
		return generator.NextID()
	}
}

func (d *Database) Close() error {
	var errString string
	for itemIdx := range d.items {
		itemErr := d.items[itemIdx].Close()
		if itemErr != nil {
			errString += fmt.Sprintf("%s error : %s", itemIdx, itemErr.Error())
		}
	}
	if errString == "" {
		return nil
	}
	return errors.New(errString)
}
