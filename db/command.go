package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// Commands  a database command session
type Commands struct {
	Querys
	kvPair map[string]interface{}
	op     op
}

func (c *Commands) NewID(tag string) (int64, error) {
	return c.ctx.database.NewID(c.ctx.ctx, c.tableName, tag)
}

func (c *Commands) prepareSQL() error {
	var (
		sqlStr   string = ""
		hasWhere        = false
		kvFlag          = false
	)
	// prepare sql script
	if c.hasRawSql == true {
		for idx := range c.conditions {
			if ok := c.checkVarsInCondition(c.conditions[idx]); ok == false {
				continue
			}
			if hasWhere == false {
				sqlStr += " WHERE "
				sqlStr += c.conditions[idx].cmd
				hasWhere = true
				continue
			}
			/* if idx == 0 {
				sqlStr += " WHERE "
				sqlStr += q.conditions[idx].cmd
				continue
			} */

			if c.conditions[idx].op == where {
				sqlStr += fmt.Sprintf(" AND %s", c.conditions[idx].cmd)
			}

			if c.conditions[idx].op == or {
				sqlStr += fmt.Sprintf(" OR %s", c.conditions[idx].cmd)
			}
		}
		c.sql += sqlStr
		return nil
	}
	if c.tableName == "" {
		return ErrInvalidCommand //errLackTableName()
	}
	if c.op == insert {
		keys := make([]string, 0, len(c.kvPair))

		for k := range c.kvPair {
			keys = append(keys, k)
		}

		sqlStr += "INSERT INTO `" + c.tableName + "` (`" + strings.Join(keys, "`"+concatDelimter+"`") + "`)  VALUES ("

		for _, key := range keys {
			sqlStr += "{" + key + "},"
		}

		sqlStr = sqlStr[:len(sqlStr)-1]
		sqlStr += ");"
	} else if c.op == delete {
		sqlStr += "DELETE FROM `" + c.tableName + "`"
		for idx := range c.conditions {
			if idx == 0 && c.conditions[idx].op == or {
				op := c.conditions[idx]
				copy(c.conditions[idx:], c.conditions[idx+1:])
				c.conditions = c.conditions[:len(c.conditions)-1]
				c.conditions = append(c.conditions, op)
			}

			/* 			for varsIdx := range c.conditions[idx].vars {
				if _, ok := c.vars[c.conditions[idx].vars[varsIdx]]; !ok {
					return ErrInvalidCommand //errLackVariable(q.conditions[idx].vars[varsIdx])
				}
			} */
		}

		for idx := range c.conditions {
			if ok := c.checkVarsInCondition(c.conditions[idx]); ok == false {
				continue
			}
			if hasWhere == false {
				sqlStr += " WHERE "
				sqlStr += c.conditions[idx].cmd
				hasWhere = true
				continue
			}

			if c.conditions[idx].op == where {
				sqlStr += fmt.Sprintf(" AND %s", c.conditions[idx].cmd)
			}

			if c.conditions[idx].op == or {
				sqlStr += fmt.Sprintf(" OR %s", c.conditions[idx].cmd)
			}
		}
	} else if c.op == update {
		sqlStr += " UPDATE `" + c.tableName + "`"
		for k := range c.kvPair {
			if kvFlag == false {
				sqlStr += " SET `" + k + "` = {" + k + "},"
				kvFlag = true
			} else {
				sqlStr += " `" + k + "` = {" + k + "},"
			}
		}
		sqlStr = sqlStr[:len(sqlStr)-1]

		for idx := range c.conditions {
			if idx == 0 && c.conditions[idx].op == or {
				op := c.conditions[idx]
				copy(c.conditions[idx:], c.conditions[idx+1:])
				c.conditions = c.conditions[:len(c.conditions)-1]
				c.conditions = append(c.conditions, op)
			}
		}

		var hasWhere = false
		for idx := range c.conditions {
			if ok := c.checkVarsInCondition(c.conditions[idx]); ok == false {
				continue
			}
			if hasWhere == false {
				sqlStr += " WHERE "
				sqlStr += c.conditions[idx].cmd
				hasWhere = true
				continue
			}

			if c.conditions[idx].op == where {
				sqlStr += fmt.Sprintf(" AND %s", c.conditions[idx].cmd)
			}

			if c.conditions[idx].op == or {
				sqlStr += fmt.Sprintf(" OR %s", c.conditions[idx].cmd)
			}
		}

		sqlStr += ";"
	}
	c.sql = sqlStr
	return nil
}

//RawSQL update command text with raw sql
func (c *Commands) RawSQL(raw string) *Commands { //, bizName string) *Querys {
	// q.bizName = bizName
	c.sql = raw
	c.hasRawSql = true
	//RawSQL更改後，Where置空， 保留Vars
	c.conditions = make([]operator, 0)
	return c
}

// For assign specific table and business for this command
func (c *Commands) For(tableName string) *Commands {
	c.tableName = tableName
	// c.bizName = bizName
	return c
}

// Update assign column and value for set expression
func (c *Commands) Update(col string, value interface{}) *Commands {
	v := RoundVar(value)
	c.kvPair[col] = v
	c.vars[col] = v
	c.op = update
	return c
}

// Updates assign args with key-value for set expresiion
func (c *Commands) Updates(args map[string]interface{}) *Commands {
	for k, v := range args {
		rv := RoundVar(v)
		c.kvPair[k] = rv
		c.vars[k] = rv
	}
	c.op = update
	return c
}

// Delete the specific rows
func (c *Commands) Delete() *Commands {
	c.op = delete
	return c
}

// Insert data into DB
func (c *Commands) Insert(col string, value interface{}) *Commands {
	v := RoundVar(value)
	c.kvPair[col] = v
	c.vars[col] = v
	c.op = insert
	return c
}

func (q *Commands) Var(name string, value interface{}) *Commands {
	q.vars[name] = RoundVar(value)
	return q
}

// Vars assign variables with key-value format to the command
func (q *Commands) Vars(vars map[string]interface{}) *Commands {
	for k := range vars {
		q.vars[k] = RoundVar(vars[k])
	}

	return q
}

// Where assign
func (q *Commands) Where(cmd string, vars ...string) *Commands {
	q.conditions = append(q.conditions, operator{op: where, cmd: cmd, vars: vars})
	return q
}

/* EX : query.OR("id ={id} ", "id")
use vars to check the var existed in the var map
*/
func (q *Commands) OR(cmd string, vars ...string) *Commands {
	q.conditions = append(q.conditions, operator{op: or, cmd: cmd, vars: vars})
	return q
}

func (c *Commands) Exec() (sql.Result, error) {
	var stmt *sql.Stmt
	var err error
	err = c.prepareSQL()

	if err != nil {
		return nil, err
	}

	defer func() {
		if stmt != nil {
			stmt.Close()
		}
		c.sql = ""
	}()

	stmt, c.sql, c.preparedArgs, err = c.ctx.PrepareContext(context.Background(), c.sql, c.vars, c.tableName)
	if err != nil {
		return nil, err
	}

	results, err := stmt.ExecContext(c.ctx.ctx, c.preparedArgs...)

	if err != nil {
		return nil, err
	}

	return results, nil

}
