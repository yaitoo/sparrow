package db

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

const (
	tagName          = "sql"
	concatDelimter   = ", "
	leftParenthesis  = "("
	rightParenthesis = ")"
)

type op int

const (
	where op = iota
	or
	insert
	delete
	update
	in
	or_in
)

type order int

const (
	asc order = iota
	desc
)

type orders struct {
	order
	cols []string
}

type operator struct {
	op   op
	cmd  string
	vars []string
}

// type condition map[string][]string

//Query a datbase query session
type Querys struct {
	vars      map[string]interface{}
	tableName string
	bizName   string
	//columns   []Column
	columns      []string
	targetModel  interface{}
	conditions   []operator
	orders       orders
	groupBy      []string
	limit        int
	offset       int
	sql          string
	preparedArgs []interface{}
	hasRawSql    bool
	// wherePos    int
	// Context     context.Context
	ctx *Context
}

//--------------------------------
func (q *Querys) For(tableName string) *Querys {
	q.tableName = tableName
	// q.bizName = bizName
	//RawSQL更改後，Where置空， 保留Vars
	q.conditions = make([]operator, 0)

	return q
}

func (q *Querys) ForModel(obj interface{}) (*Querys, error) {
	v := reflect.ValueOf(obj)
	if v.Elem().Type().Kind() != reflect.Struct {
		return nil, ErrInvalidObject
	}

	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		q.tableName = t.Elem().Name()
	} else {
		q.tableName = t.Name()
	}

	//RawSQL更改後，Where置空， 保留Vars
	q.conditions = make([]operator, 0)

	return q, nil
}

func (q *Querys) NewID(tag string) (int64, error) {
	if q.tableName == "" || tag == "" {
		return 0, nil
	}
	// check table and tag exist in the latest rule

	if tableGen, ok := Shardings[q.tableName]; ok {
		if tagGen, ok := tableGen[tag]; ok {
			shardingId, err := tagGen.NextID()
			if err != nil {
				return 0, ErrClockMovedBackwards
			}
			return shardingId, nil
		}
	}
	return 0, nil
}

//---------------------------------
// func (q *Querys) Select(columns ...Column) *Querys {
func (q *Querys) Select(columns ...string) *Querys {
	q.columns = columns
	return q
}

func (q *Querys) SelectModel(object interface{}) *Querys {
	q.columns = []string{}

	q.targetModel = object
	return q
}

/* EX : query.Where("id = {id}", "id")
use vars to check the var existed in the var map
*/
func (q *Querys) Where(cmd string, varNames ...string) *Querys {
	q.conditions = append(q.conditions, operator{op: where, cmd: cmd, vars: varNames})
	return q
}

/* EX : query.WhereOrIN("id in {id}", "id")
該欄位不能是參與分片的欄位
*/
func (q *Querys) WhereOrIN(cmd string, varNames ...string) *Querys {
	q.conditions = append(q.conditions, operator{op: or_in, cmd: cmd, vars: varNames})
	return q
}

/* EX : query.WhereIN("id in {id}", "id")
該欄位不能是參與分片的欄位
varNames 指定关联变量名称，如果有任一关联变量名称没有赋值，则这段cmd会被忽略，不参与最终SQL
*/
func (q *Querys) WhereIN(cmd string, varNames ...string) *Querys {
	q.conditions = append(q.conditions, operator{op: in, cmd: cmd, vars: varNames})
	return q
}

/* EX : query.WhereOR("id ={id} ", "id")
use vars to check the var existed in the var map
*/
func (q *Querys) WhereOR(cmd string, vars ...string) *Querys {
	q.conditions = append(q.conditions, operator{op: or, cmd: cmd, vars: vars})
	return q
}

/* EX : query.IN("id in (a,b,c) ", "id")
use vars to check the var existed in the var map
*/
/* func (q *Querys) IN(cmd string, vars ...[]string) *Querys {
	q.search.IN(cmd, vars)
	return q
} */

//VarIN 添加IN对应的变量，值为数组
func (q *Querys) VarIN(name string, values ...interface{}) *Querys {
	if len(values) > 0 {
		q.vars[name] = values
	}

	return q
}

//Var 添加变量
func (q *Querys) Var(name string, value interface{}) *Querys {

	q.vars[name] = RoundVar(value)
	return q
}

//Vars 批量添加变量， 同名覆盖
func (q *Querys) Vars(vars map[string]interface{}) *Querys {
	for k, v := range vars {
		q.vars[k] = RoundVar(v)
	}

	return q
}

func (q *Querys) OrderBy(columns ...string) *Querys {
	q.orders.order = asc
	q.orders.cols = columns

	return q
}

func (q *Querys) OrderByDescending(columns ...string) *Querys {
	q.orders.order = desc
	q.orders.cols = columns
	return q
}

func (q *Querys) RawSQL(raw string) *Querys { //, bizName string) *Querys {
	// q.bizName = bizName
	q.sql = raw
	q.hasRawSql = true
	//RawSQL更改後，Where置空， 保留Vars
	q.conditions = make([]operator, 0)
	return q
}

func (q *Querys) GroupBy(columns ...string) *Querys {
	q.groupBy = columns
	return q
}

//Take 設定分頁頁碼和每頁大小
func (q *Querys) Take(pageIndex, pageSize int) *Querys {
	if pageSize > 0 {
		q.limit = pageSize
	} else {
		q.limit = 0
	}

	if pageIndex > 1 {
		q.offset = (pageIndex - 1) * pageSize
	} else {
		q.offset = 0
	}
	return q
}

func (q *Querys) checkVarsInCondition(operator operator) bool {
	for idx := range operator.vars {
		if _, ok := q.vars[operator.vars[idx]]; ok == false {
			return false
		}
	}
	return true
}

func getInSqlStr(vars map[string]interface{}, condition operator, needOp bool) (string, error) {
	var (
		tempCmd        string
		charIndex      int
		rightCharIndex int
		varArr         []string
	)
	if charIndex = strings.Index(condition.cmd, "{"); charIndex == -1 {
		return "", errors.Wrap(ErrInvalidCommand, condition.cmd)
	}

	if rightCharIndex = strings.LastIndex(condition.cmd, "}"); rightCharIndex == -1 {
		return "", errors.Wrap(ErrInvalidCommand, condition.cmd)
	}

	varName := condition.cmd[charIndex+1 : rightCharIndex]
	values, ok := vars[varName]
	if !ok {
		return "", errors.Wrap(ErrInvalidCommand, "["+varName+"] is missing for ["+condition.cmd+"]")
	}

	valueAarray := values.([]interface{})

	tempCmd = condition.cmd[:charIndex] + "("
	for varIdx := range valueAarray {
		col := fmt.Sprintf("%s_in_%d", varName, varIdx)
		varArr = append(varArr, fmt.Sprintf("{%s}", col))
		vars[col] = valueAarray[varIdx]
	}

	if needOp == false {
		return fmt.Sprintf(" %s %s)", tempCmd, strings.Join(varArr, concatDelimter)), nil
	}
	if condition.op == in {
		return fmt.Sprintf(" AND %s %s)", tempCmd, strings.Join(varArr, concatDelimter)), nil
	}
	if condition.op == or_in {
		return fmt.Sprintf(" OR %s %s)", tempCmd, strings.Join(varArr, concatDelimter)), nil
	}
	return "", ErrInvalidCommand
}

func (q *Querys) prepareSQL() error {
	var sqlStr string = ""
	if q.hasRawSql == false {
		if q.tableName == "" {
			return ErrInvalidCommand
		}

		if q.targetModel != nil {
			if columns := getSQLColumns(q.targetModel); columns != nil {
				q.columns = columns
			}
		}
		// prepare sql script
		sqlStr += "SELECT "
		// sqlStr += strings.Join(q.columns, concatDelimter)

		cols := make([]string, len(q.columns))
		for i := 0; i < len(q.columns); i++ {
			cols[i] = fmt.Sprintf("`%s`.", q.tableName) + fmt.Sprintf("`%s`", q.columns[i])
		}
		sqlStr += strings.Join(cols, concatDelimter)

		sqlStr += fmt.Sprintf(" FROM `%s`", q.tableName)
	}

	for idx := range q.conditions {
		if idx == 0 && (q.conditions[idx].op == or || q.conditions[idx].op == or_in) {
			op := q.conditions[idx]
			copy(q.conditions[idx:], q.conditions[idx+1:])
			q.conditions = q.conditions[:len(q.conditions)-1]
			q.conditions = append(q.conditions, op)
		}
		break
	}
	var hasWhere = false
	for idx := range q.conditions {
		if ok := q.checkVarsInCondition(q.conditions[idx]); ok == false {
			continue
		}
		if hasWhere == false {
			sqlStr += " WHERE "
			if q.conditions[idx].op == in {
				/* var (
					tempCmd   string
					charIndex int
					varArr    []string
				)
				if charIndex = strings.Index(q.conditions[idx].cmd, leftParenthesis); charIndex == -1 {
					return ErrInvalidCommand
				}
				tempCmd = q.conditions[idx].cmd[:charIndex+1]
				for varIdx := range q.vars[q.conditions[idx].vars[0]].([]interface{}) {
					col := fmt.Sprintf("%s_in_%d", q.conditions[idx].vars[0], varIdx)
					varArr = append(varArr, fmt.Sprintf("{%s}", col))
					q.vars[col] = q.vars[q.conditions[idx].vars[0]].([]interface{})[varIdx]
				}

				sqlStr += fmt.Sprintf(" %s %s %s ", tempCmd, strings.Join(varArr, concatDelimter), rightParenthesis) */
				inSql, err := getInSqlStr(q.vars, q.conditions[idx], false)
				if err != nil {
					return err
				}
				sqlStr += inSql
			} else {
				sqlStr += q.conditions[idx].cmd
			}
			hasWhere = true
			continue
		}

		if q.conditions[idx].op == where {
			sqlStr += fmt.Sprintf(" AND %s", q.conditions[idx].cmd)
		}

		if q.conditions[idx].op == or {
			sqlStr += fmt.Sprintf(" OR %s", q.conditions[idx].cmd)
		}

		if q.conditions[idx].op == in || q.conditions[idx].op == or_in {
			/* var (
				tempCmd   string
				charIndex int
				varArr    []string
			)
			if charIndex = strings.Index(q.conditions[idx].cmd, leftParenthesis); charIndex == -1 {
				return ErrInvalidCommand
			}
			tempCmd = q.conditions[idx].cmd[:charIndex+1]
			for varIdx := range q.vars[q.conditions[idx].vars[0]].([]interface{}) {
				col := fmt.Sprintf("%s_in_%d", q.conditions[idx].vars[0], varIdx)
				varArr = append(varArr, fmt.Sprintf("{%s}", col))
				q.vars[col] = q.vars[q.conditions[idx].vars[0]].([]interface{})[varIdx]
			}

			sqlStr += fmt.Sprintf(" AND %s %s %s ", tempCmd, strings.Join(varArr, concatDelimter), rightParenthesis) */
			inSql, err := getInSqlStr(q.vars, q.conditions[idx], true)
			if err != nil {
				return err
			}
			sqlStr += inSql
		}

		/* 	if q.conditions[idx].op == or_in {
			var (
				tempCmd   string
				charIndex int
				varArr    []string
			)
			if charIndex = strings.Index(q.conditions[idx].cmd, leftParenthesis); charIndex == -1 {
				return ErrInvalidCommand
			}
			tempCmd = q.conditions[idx].cmd[:charIndex+1]
			for varIdx := range q.vars[q.conditions[idx].vars[0]].([]interface{}) {
				col := fmt.Sprintf("%s_in_%d", q.conditions[idx].vars[0], varIdx)
				varArr = append(varArr, fmt.Sprintf("{%s}", col))
				q.vars[col] = q.vars[q.conditions[idx].vars[0]].([]interface{})[varIdx]
			}

			sqlStr += fmt.Sprintf(" OR %s %s %s ", tempCmd, strings.Join(varArr, concatDelimter), rightParenthesis)
		} */
	}

	if len(q.orders.cols) > 0 {
		sqlStr += " ORDER BY "
		sqlStr += strings.Join(q.orders.cols, concatDelimter)
		if q.orders.order == asc {
			sqlStr += " ASC"
		} else {
			sqlStr += " DESC"
		}
	}

	if len(q.groupBy) > 0 {
		sqlStr += " GROUP BY "
		sqlStr += strings.Join(q.groupBy, concatDelimter)
	}

	if q.offset > 0 && q.limit > 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d, %d", q.offset, q.limit)
	} else if q.offset == 0 && q.limit > 0 {
		sqlStr += fmt.Sprintf(" LIMIT %d", q.limit)
	}

	q.sql += sqlStr

	return nil
}

//First 获取排序第一的记录，填入指定的对象
func (q *Querys) First(obj interface{}) error {
	var stmt *sql.Stmt
	var err error

	q.limit = 1
	err = q.prepareSQL()

	if err != nil {
		return err
	}

	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	stmt, q.sql, q.preparedArgs, err = q.ctx.PrepareContext(context.Background(), q.sql, q.vars, q.tableName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	rows, err := q.queryRowsContext(stmt)

	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		rows.Close()
		q.limit = 0
		q.offset = 0
		q.sql = ""
		q.hasRawSql = false
	}()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	objKind := reflect.ValueOf(obj)
	if objKind.Kind() == reflect.Ptr {
		if objKind.Elem().Kind() == reflect.Struct {

			err = ToStruct(q.ctx.ctx, rows, obj)
			if err != nil {
				logger.Error(err, " :", q.sql)
				//return ErrInvalidObject //err
				return err
			}
			return nil

		}

		err = rows.Scan(obj)
		if err != nil {
			logger.Error(err)
			return err
		}
		return nil

	}

	// return value
	q.sql = ""
	return nil
}

//Count 计算返回结果
func (q *Querys) Count() (int, error) {
	q.columns = []string{"count(*)"}

	var stmt *sql.Stmt
	var err error

	err = q.prepareSQL()

	if err != nil {
		return 0, err
	}

	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	stmt, q.sql, q.preparedArgs, err = q.ctx.PrepareContext(context.Background(), q.sql, q.vars, q.tableName)
	if err != nil {
		return 0, err
	}

	rows, err := q.queryRowsContext(stmt)

	if err != nil {
		return 0, err
	}

	defer func() {
		rows.Close()
		q.limit = 0
		q.offset = 0
		q.sql = ""
		q.hasRawSql = false
	}()
	if err != nil {
		return 0, err
	}

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}

//Find 把查询回的数据填入指定的对象，必须使用引用对象，否则会返回ErrInvalidObject异常
func (q *Querys) Find(objects interface{}) error {

	var stmt *sql.Stmt
	var err error

	err = q.prepareSQL()

	if err != nil {
		return err
	}

	defer func() {
		if stmt != nil {
			stmt.Close()
		}
	}()

	stmt, q.sql, q.preparedArgs, err = q.ctx.PrepareContext(context.Background(), q.sql, q.vars, q.tableName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	rows, err := q.queryRowsContext(stmt)

	if err != nil {
		logger.Error(err)
		return err
	}

	defer func() {
		rows.Close()
		q.limit = 0
		q.offset = 0
		q.sql = ""
		q.hasRawSql = false
	}()
	if err != nil {
		logger.Error(err)
		return err
	}
	// return ToStruct(rows, objects)
	objKind := reflect.ValueOf(objects)
	if objKind.Kind() == reflect.Ptr {
		if objKind.Elem().Kind() == reflect.Struct { //非slice，需要处理ErrNoRows
			if rows.Next() {
				err = ToStruct(q.ctx.ctx, rows, objects)
				if err != nil {
					logger.Error(err)
					//return ErrInvalidObject //err
					return err
				}

				return nil
			}

			return sql.ErrNoRows

		} else if objKind.Elem().Kind() == reflect.Slice { //集合，忽略ErrNoRows
			err = ScanAll(q.ctx.ctx, rows, objects, false)
			if err != nil {
				logger.Error(err)
				//return ErrInvalidObject //err
				return err
			}
			return nil
		} else {
			if rows.Next() {
				err = rows.Scan(objects)
				if err != nil {
					logger.Error(err)
					//return ErrInvalidObject // err
					return err
				}

				return nil
			}
			return sql.ErrNoRows
		}

	}

	return ErrInvalidObject
}

//QueryRow 提供底层的QueryRow操作，暴露底层sql操作对象代替ORM工作
func (q *Querys) QueryRow(scan func(row *sql.Row) error) error {
	var stmt *sql.Stmt
	var err error

	err = q.prepareSQL()

	if err != nil {
		return err
	}

	stmt, q.sql, q.preparedArgs, err = q.ctx.PrepareContext(context.Background(), q.sql, q.vars, q.tableName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	row := q.queryRowContext(stmt)
	defer func() {
		q.limit = 0
		q.offset = 0
		q.sql = ""
		q.hasRawSql = false
	}()
	return scan(row)

}

//Query 提供底层的操作接口，暴露sql的操作对象代替ORM工作
func (q *Querys) Query(scan func(rows *sql.Rows) error) error {
	var stmt *sql.Stmt
	var err error

	err = q.prepareSQL()

	if err != nil {
		return err
	}

	stmt, q.sql, q.preparedArgs, err = q.ctx.PrepareContext(context.Background(), q.sql, q.vars, q.tableName)
	if err != nil {
		logger.Errorln(err)
		return err
	}

	rows, err := q.queryRowsContext(stmt)
	defer func() {
		rows.Close()
		q.limit = 0
		q.offset = 0
		//	q.sql = ""
		q.hasRawSql = false
	}()

	if err != nil {
		logger.Error(err)
		return err
	}

	return scan(rows)
}

/* -------------------------------- */
func (s *Querys) Limit(limit int) *Querys {
	s.limit = limit
	return s
}

func (s *Querys) Offset(offset int) *Querys {
	s.offset = offset
	return s
}

/* -------------------- */
func (s *Querys) queryRow(db *sql.DB) *sql.Row {
	return db.QueryRowContext(s.ctx.ctx, s.sql)
}

func (s *Querys) queryRows(db *sql.DB) (*sql.Rows, error) {
	return db.QueryContext(s.ctx.ctx, s.sql)
}

func (s *Querys) queryRowContext(stmt *sql.Stmt) *sql.Row {
	return stmt.QueryRowContext(s.ctx.ctx, s.preparedArgs...)
}

func (s *Querys) queryRowsContext(stmt *sql.Stmt) (*sql.Rows, error) {
	return stmt.QueryContext(s.ctx.ctx, s.preparedArgs...)
}

func (s *Querys) String() string {
	return s.sql
}
