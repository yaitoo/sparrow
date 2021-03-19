package db

import (
	"database/sql"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/types"
)

//Command a db command
type Command struct {
	SQLBuilder
	ctx *Context
}

//Exec execute command in db
func (c *Command) Exec(shardingTime *time.Time) (sql.Result, error) {
	return c.ExecSQL(c.String(), shardingTime)
}

//ExecSQL execute command in db
func (c *Command) ExecSQL(s string, shardingTime *time.Time) (sql.Result, error) {
	if shardingTime != nil {
		shardings := getShardings(s, shardingTime, nil)

		if len(shardings) > 0 {
			s = strings.Replace(s, shardings[0].token, shardings[0].value, -1)
		}
	}

	startTime := time.Now()
	stmt, args, formattedCMD, err := c.ctx.prepare(s, c.vars)
	if stmt != nil {
		defer stmt.Close()

		r, err := stmt.Exec(args...)

		c.ctx.logSlowSQL(formattedCMD, time.Since(startTime))

		if err != nil {
			logger.Warnln(err, formattedCMD)
			return nil, err
		}

		return r, nil

	}

	logger.Warnln(err, formattedCMD)
	return nil, err
}

//First return first columan as a single value
func (c *Command) First(cmd string, shardingTime *time.Time) string {
	val := &types.String{}
	s := cmd + c.buildWhere()

	if err := c.ctx.findWith(val, s, c.vars); err != nil {
		logger.Warnln(err)
	}

	return val.GetValue()
}
