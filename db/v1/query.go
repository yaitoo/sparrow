package db

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/types"
)

//Query query data in db
type Query struct {
	SQLBuilder
	ctx *Context
}

//Find result single value with sql in db
func (q *Query) Find(model interface{}, shardingTime *time.Time) error {

	s := q.String()

	if shardingTime != nil {
		shardings := getShardings(s, shardingTime, nil)

		if len(shardings) > 0 {
			s = strings.Replace(s, shardings[0].token, shardings[0].value, -1)
		}
	}

	return q.ctx.findWith(model, s, q.vars)

}

//FindAll return all values with sql in db
func (q *Query) FindAll(list interface{}, shardingTimeFrom, shardingTimeTo *time.Time) error {

	listType := reflect.TypeOf(list)

	if listType.Kind() == reflect.Ptr && listType.Elem().Kind() == reflect.Slice {

	} else {
		return fmt.Errorf("type: %s is unspported yet", listType.Elem().Kind())
	}

	results := reflect.ValueOf(list).Elem()
	sliceType := listType.Elem()

	elementType := sliceType.Elem()
	isPtr := elementType.Kind() == reflect.Ptr
	if isPtr {
		elementType = elementType.Elem()
	}

	results = reflect.ValueOf(list).Elem()

	a := loadAnnotationByType(elementType)

	cmd := q.String()

	shardings := getShardings(cmd, shardingTimeFrom, shardingTimeTo)

	resultChan := make(chan *queryResult, len(shardings))

	for _, it := range shardings {

		go func(it sharding) {

			startTime := time.Now()
			items, err := q.execFindAll(strings.Replace(cmd, it.token, it.value, -1), a, elementType)

			resultChan <- &queryResult{items: items, err: err, queryDuration: time.Since(startTime)}

		}(it)
	}

	var totalDuration time.Duration
	var err error
	for i := 0; i < len(shardings); i++ {

		result := <-resultChan

		if result != nil {

			totalDuration += result.queryDuration

			if err != nil {
				continue
			}
			if result.err != nil {
				err = result.err
				continue
			}

			for _, it := range result.items {
				if isPtr {
					results.Set(reflect.Append(results, reflect.ValueOf(it)))
				} else {
					results.Set(reflect.Append(results, reflect.ValueOf(it).Elem()))
				}
			}
		}

	}

	q.ctx.logSlowSQL(cmd, totalDuration)

	return err

}

func (q *Query) execFindAll(cmd string, a *annotation, elementType reflect.Type) ([]interface{}, error) {
	stmt, args, fmtSQL, err := q.ctx.prepare(cmd, q.vars)

	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, 100)

	if stmt != nil {
		defer stmt.Close()

		rows, err := stmt.Query(args...)
		defer rows.Close()

		if err != nil {
			logger.Warnln(err, fmtSQL)
			return nil, err
		}

		cols, err := rows.Columns()
		if err != nil {
			logger.Warnln(err, fmtSQL)
			return nil, err
		}

		for rows.Next() {

			item := reflect.New(elementType).Interface()

			fglobalizer, ok := item.(types.Contexter)

			if ok {
				fglobalizer.SetContext(q.ctx.ctx)
			}

			if _, ok := item.(sql.Scanner); ok {
				if err := rows.Scan(item); err != nil {
					logger.Warnln(err, fmtSQL)
					return nil, err
				}
			} else {
				if a.isPrimitive {
					if err := rows.Scan(item); err != nil {
						logger.Warnln(err, fmtSQL)
						return nil, err
					}
				} else {
					if err := rows.Scan(a.GetScanPtrs(q.ctx, item, cols)...); err != nil {
						logger.Warnln(err, fmtSQL)
						return nil, err
					}
				}

			}

			items = append(items, item)

		}

		if err := rows.Close(); err != nil {
			return nil, err
		}

		return items, nil
	}

	return nil, nil
}

// func (q *Query) ScanRow(callback func(*sql.Row) error) error {
// 	cmd := q.String()

// 	stmt, args, err := q.Context.Prepare(cmd, q.values)

// 	if stmt != nil {
// 		defer stmt.Close()
// 		row := stmt.QueryRow(args...)

// 		if err := callback(row); err != nil {
// 			return err
// 		} else {
// 			return nil
// 		}

// 	} else {
// 		return err
// 	}
// }

// func (q *Query) ScanRows(callback func(*sql.Rows) error) error {

// 	cmd := q.String()

// 	stmt, args, err := q.Context.Prepare(cmd, q.values)

// 	if stmt != nil {
// 		defer stmt.Close()
// 		rows, err := stmt.Query(args...)
// 		if err != nil {
// 			return err
// 		}
// 		defer rows.Close()
// 		for rows.Next() {
// 			if err := callback(rows); err != nil {
// 				tracer.Error("db", cmd, err)
// 				return err
// 			}
// 		}

// 		if err != nil {
// 			tracer.Error("db", cmd, err)
// 			return err

// 		}

// 		return nil
// 	} else {
// 		tracer.Error("db", cmd, err)
// 		return err
// 	}
// }

type queryResult struct {
	items         []interface{}
	err           error
	queryDuration time.Duration
}
