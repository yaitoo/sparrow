package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/yaitoo/sparrow/log"
	"github.com/yaitoo/sparrow/types"
)

//Context db context
type Context struct {
	sync.Mutex
	ctx             context.Context
	db              *sql.DB
	tx              *sql.Tx
	slowSQLDuration *time.Duration
}

var (
	regexSQLVarToken = regexp.MustCompile(`{(.)+?}`)

	logger = log.NewLogger("db")
)

//NewContext create a db context instance, is not thread safe
func NewContext(ctx context.Context, driverName, dataSourceName string) *Context {
	dc := &Context{}
	dc.ctx = ctx

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		logger.Error(err)
		return nil
	}

	dc.db = db

	return dc
}

//Clone clone a new db context with new context.Context
func (c *Context) Clone(ctx context.Context) *Context {
	dc := &Context{}
	dc.db = c.db
	if ctx == nil {
		dc.ctx = c.ctx
	} else {
		dc.ctx = ctx
	}

	dc.slowSQLDuration = c.slowSQLDuration

	return dc
}

//RawDB return raw sql.DB instance
func (c *Context) RawDB() *sql.DB {
	return c.db
}

//Begin start a db transaction
func (c *Context) Begin() error {
	c.Lock()
	defer c.Unlock()

	if c.tx != nil {
		err := c.Rollback()
		c.tx = nil

		if err != nil {
			return err
		}
	}

	tx, err := c.db.Begin()

	if err != nil {
		return err
	}

	c.tx = tx
	return nil
}

//Commit commit the db transaction if it is valid
func (c *Context) Commit() error {
	c.Lock()
	defer c.Unlock()
	if c.tx != nil {
		err := c.tx.Commit()
		if err != nil {
			c.tx.Rollback()
		}
		c.tx = nil
		return err
	}

	return nil
}

//Rollback rollback the db transaction if it is valid
func (c *Context) Rollback() error {
	c.Lock()
	defer c.Unlock()
	if c.tx != nil {
		err := c.tx.Rollback()
		c.tx = nil
		return err
	}

	return nil
}

//TraceSlowSQL turn on/off slow sql tracing
func (c *Context) TraceSlowSQL(status bool, t *time.Duration) {
	c.Lock()
	defer c.Unlock()
	if status {
		c.slowSQLDuration = t
	} else {
		c.slowSQLDuration = nil
	}
}

//NewCommand return a new command
func (c *Context) NewCommand() *Command {
	return &Command{
		ctx: c,
	}
}

//NewQuery return a new query
func (c *Context) NewQuery() *Query {
	return &Query{
		ctx: c,
	}
}

//Prepare creates a prepared statement for use within a transaction.
func (c *Context) prepare(cmd string, vars map[string]interface{}) (*sql.Stmt, []interface{}, string, error) {
	tokens := regexSQLVarToken.FindAllStringIndex(cmd, -1)
	n := len(tokens)

	var stmt *sql.Stmt
	var err error

	if n > 0 {
		args := make([]interface{}, n, n)

		s := ""
		formattedCMD := ""
		i := 0
		for k, v := range tokens {
			s += cmd[i:v[0]] + "?"

			name := cmd[v[0]+1 : v[1]-1]

			val, ok := vars[strings.ToLower(name)]
			if ok {
				args[k] = val

				formattedCMD += cmd[i:v[0]] + fmt.Sprintf("%s", val)
			} else {
				return nil, nil, cmd, fmt.Errorf("db: {%s} is missing in %s", name, cmd)
			}
			i = v[1]
		}

		if i < len(cmd) {
			s += cmd[i:]
			formattedCMD += cmd[i:]
		}

		if c.tx != nil {
			stmt, err = c.tx.Prepare(s)
		} else {
			stmt, err = c.db.Prepare(s)
		}

		return stmt, args, formattedCMD, err

	}

	if c.tx != nil {
		stmt, err = c.tx.Prepare(cmd)
	} else {
		stmt, err = c.db.Prepare(cmd)
	}

	return stmt, nil, cmd, err

}

func (c *Context) findWith(model interface{}, cmd string, values map[string]interface{}) error {

	startTime := time.Now()
	stmt, args, fmtSQL, err := c.prepare(cmd, values)
	if stmt != nil {
		defer stmt.Close()

		r, err := stmt.Query(args...)

		c.logSlowSQL(fmtSQL, time.Since(startTime))

		if err != nil {
			logger.Warnln(err, fmtSQL)
			return err
		}
		defer r.Close()

		var dest []interface{}

		if _, ok := model.(sql.Scanner); ok {
			dest = make([]interface{}, 1, 1)
			contexter, ok := model.(types.Contexter)
			if ok {
				contexter.SetContext(c.ctx)
			}

			dest[0] = model
		} else {

			t := reflect.TypeOf(model)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}

			if t.Kind() == reflect.Struct {
				a := loadAnnotation(model)

				cols, err := r.Columns()
				if err != nil {
					logger.Warnln(err, fmtSQL)
					return err
				}

				dest = a.GetScanPtrs(c, model, cols)
			} else {
				dest = make([]interface{}, 1, 1)
				dest[0] = model
			}

		}

		for _, dp := range dest {
			if _, ok := dp.(*sql.RawBytes); ok {
				return errors.New("sql: RawBytes isn't allowed on Row.Scan")
			}
		}

		if !r.Next() {
			if err := r.Err(); err != nil {
				return err
			}
			return sql.ErrNoRows
		}
		err = r.Scan(dest...)
		if err != nil {
			return err
		}
		// Make sure the query can be processed to completion with no errors.
		if err := r.Close(); err != nil {
			return err
		}

		return nil

	}

	return err

}

func (c *Context) logSlowSQL(sql string, duration time.Duration) {
	slowSQLDuration := c.slowSQLDuration
	if slowSQLDuration == nil {
		return
	}

	if *slowSQLDuration <= duration {
		logger.Warnln("SlowSQL ", duration, " : ", sql)
	}
}
