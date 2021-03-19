package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/linq"
	"github.com/yaitoo/sparrow/types"
)

//SQLBuilder  sql builder
type SQLBuilder struct {
	vars         map[string]interface{}
	wheres       []whereClause
	orderColumns map[string]string
	orderbys     []string
	limit        string
	sql          string
}

type whereClause struct {
	op     string
	clause string
}

//RawSQL set raw sql
func (sb *SQLBuilder) RawSQL(sql string) {
	sb.sql = sql
}

//Var add/update variable
func (sb *SQLBuilder) Var(name string, value interface{}, predicate bool) {
	if predicate == false {
		return
	}

	if sb.vars == nil {
		sb.vars = make(map[string]interface{})
	}

	if types.IsEmpty(name) {
		return
	}

	key := strings.ToLower(name)

	switch v := value.(type) {
	case float32:
		sb.vars[key] = types.Round32(v, .5, 2)
	case *float32:
		if v == nil {
			sb.vars[key] = nil
		} else {
			sb.vars[key] = types.Round32(*v, .5, 2)
		}
	case float64:
		sb.vars[key] = types.Round64(v, .5, 2)
	case *float64:
		if v == nil {
			sb.vars[key] = nil
		} else {
			sb.vars[key] = types.Round64(*v, .5, 2)
		}

	case time.Time:
		sb.vars[key] = v.Round(time.Microsecond)
	case *time.Time:
		if v == nil {
			sb.vars[key] = nil
		} else {
			sb.vars[key] = (*v).Round(time.Microsecond)
		}
	default:
		sb.vars[key] = value
	}

}

//Vars return a Vars
func (sb *SQLBuilder) Vars() *Vars {
	return &Vars{builder: sb}
}

//WhereAnd add a and where clause
func (sb *SQLBuilder) WhereAnd(clause string) {
	if sb.wheres == nil {
		sb.wheres = make([]whereClause, 0, 10)
	}

	if types.IsNotEmpty(clause) {
		sb.wheres = append(sb.wheres, whereClause{op: " AND ", clause: clause})
	}
}

//WhereOr add a or where clause
func (sb *SQLBuilder) WhereOr(clause string) {
	if sb.wheres == nil {
		sb.wheres = make([]whereClause, 0, 10)
	}

	if types.IsNotEmpty(clause) {
		sb.wheres = append(sb.wheres, whereClause{op: " OR ", clause: clause})
	}
}

//SetOrderColumn add orderable column mapping, prevent sql injection
func (sb *SQLBuilder) SetOrderColumn(clientColName, dbColName string) {
	if sb.orderColumns == nil {
		sb.orderColumns = make(map[string]string)
	}

	sb.orderColumns[strings.ToLower(clientColName)] = strings.ToLower(dbColName)
}

//OrderBy set order by clause
func (sb *SQLBuilder) OrderBy(clause string) {

	if sb.orderColumns == nil {
		return
	}

	if sb.orderbys == nil {
		sb.orderbys = make([]string, 0, 10)
	}

	items := linq.FromString(clause, ",").
		Where(func(i string) bool { return types.IsNotEmpty(i) }).
		ToArray()

	for _, item := range items {

		parts := strings.Split(item, " ")
		if len(parts) > 0 {
			clientColName := strings.ToLower(strings.Trim(parts[0], " "))
			direction := ""
			if len(parts) > 1 {
				direction = strings.ToLower(strings.Trim(parts[1], " "))
			}

			if !(direction == "desc" || direction == "asc") {
				direction = ""
			}

			dbColName, ok := sb.orderColumns[clientColName]
			if ok {
				sb.orderbys = append(sb.orderbys, dbColName+" "+direction)
			}

		}

	}

}

//Limit set limit clause,ignore when size is 0
func (sb *SQLBuilder) Limit(size, index int64) {
	i, s := index, size

	if s > 0 {
		if i < 1 {
			i = 1
		}

		sb.limit = fmt.Sprintf(" LIMIT %v,%v ", (i-1)*s, s)
	}
}

func (sb *SQLBuilder) buildWhere() string {
	if len(sb.wheres) > 0 {

		list := make([]string, 0, len(sb.wheres))
		for _, it := range sb.wheres {
			tokens := regexSQLVarToken.FindAllString(it.clause, -1)
			if len(tokens) > 0 {
				//token = {token}
				for _, token := range tokens {
					if _, ok := sb.vars[strings.ToLower(token[1:len(token)-1])]; !ok {
						goto skip
					}
				}

			}

			if len(list) == 0 {
				list = append(list, it.clause)
			} else {
				list = append(list, it.op+it.clause)
			}

		skip:
			continue

		}

		return fmt.Sprintln(" WHERE " + strings.Join(list, " "))
	}

	return ""
}

func (sb *SQLBuilder) buildOrderBy() string {
	if len(sb.orderbys) > 0 {
		return " ORDER BY " + strings.Join(sb.orderbys, ", ")
	}

	return ""
}

//String implements Stringer, return final sql
func (sb *SQLBuilder) String() string {

	where := sb.buildWhere()
	orderby := sb.buildOrderBy()
	sql := sb.sql

	if len(where) > 0 && strings.Index(sql, "/*where*/") > -1 {
		sql = strings.Replace(sql, "/*where*/", where, -1)
	} else {
		sql = sql + where
	}

	if len(orderby) > 0 && strings.Index(sql, "/*orderby*/") > -1 {
		sql = strings.Replace(sql, "/*orderby*/", orderby, -1)
	} else {
		sql = sql + orderby
	}

	if len(sb.limit) > 0 && strings.Index(sql, "/*limit*/") > -1 {
		sql = strings.Replace(sql, "/*limit*/", sb.limit, -1)
	} else {
		sql = sql + sb.limit
	}

	return sql
}
