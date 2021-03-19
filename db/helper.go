package db

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/serenize/snaker"
	"github.com/yaitoo/sparrow/types"
)

var (
	sqlColumnsCache = make(map[reflect.Type][]string)
	sqlColumnsMutex sync.RWMutex
)

func getSQLColumns(target interface{}) []string {

	v := reflect.Indirect(reflect.ValueOf(target))

	if v.Kind() == reflect.Struct {

		t := v.Type()

		sqlColumnsMutex.RLock()
		defer sqlColumnsMutex.RUnlock()

		columns, ok := sqlColumnsCache[t]

		if ok {
			return columns
		}

		for i := 0; i < v.NumField(); i++ {
			columnName := strings.TrimSpace(v.Type().Field(i).Tag.Get("sql"))
			if columnName != "" {
				if columnName != "-" {
					columns = append(columns, columnName)
				}
			} else {
				columns = append(columns, snaker.CamelToSnake(v.Type().Field(i).Name))
			}
		}

		go cacheSQLColumns(t, columns)

		return columns
	}

	return nil
}

func cacheSQLColumns(key reflect.Type, columns []string) {
	sqlColumnsMutex.Lock()
	defer sqlColumnsMutex.Unlock()

	sqlColumnsCache[key] = columns
}

//RoundVar 修复变量精度
func RoundVar(value interface{}) interface{} {

	if value == nil {
		return value
	}

	switch v := value.(type) {
	case float32:
		return types.Round32(v, .5, 2)
	case *float32:
		if v == nil {
			return nil
		}

		return types.Round32(*v, .5, 2)

	case float64:
		return types.Round64(v, .5, 2)
	case *float64:
		if v == nil {
			return nil
		}
		return types.Round64(*v, .5, 2)

	case time.Time:
		return v.Round(time.Microsecond)
	case *time.Time:
		if v == nil {
			return nil
		}

		return (*v).Round(time.Microsecond)

	default:
		return value
	}
}
