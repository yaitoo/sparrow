package db

import (
	"reflect"
	"strings"
	"sync"

	"github.com/yaitoo/sparrow/types"
)

type annotation struct {
	isPrimitive bool
	Columns     map[string]string
}

var (
	locker            sync.RWMutex
	cachedAnnotations = make(map[string]*annotation)
)

func loadAnnotation(m interface{}) *annotation {
	return loadAnnotationByType(reflect.TypeOf(m).Elem())
}

func loadAnnotationByType(t reflect.Type) *annotation {
	if t.PkgPath() == "" {
		return &annotation{isPrimitive: true}
	}

	aname := t.PkgPath() + "/" + t.Name()
	a, ok := cachedAnnotations[aname]

	if ok {
		return a
	}

	locker.Lock()
	defer locker.Unlock()

	a, ok = cachedAnnotations[aname]

	if ok {
		return a
	}

	a = &annotation{}
	a.Columns = make(map[string]string)

	if t.Kind() != reflect.Struct {
		return a
	}

	var fields = a.getFields(t)
	for _, field := range fields {

		col := field.Tag.Get("db")
		if types.IsNotEmpty(col) {
			a.Columns[strings.ToLower(col)] = field.Name
		} else {
			a.Columns[strings.ToLower(field.Name)] = field.Name
		}

	}

	return a
}

func (at *annotation) GetScanPtrs(dc *Context, m interface{}, columns []string) []interface{} {

	ptrs := make([]interface{}, len(columns), len(columns))

	values := reflect.ValueOf(m).Elem()

	for i, name := range columns {
		cname := strings.Replace(strings.ToLower(name), "_", "", -1)
		fname, ok := at.Columns[cname]
		if ok {
			fval := values.FieldByName(fname)

			if fval.CanAddr() {
				fref := fval.Addr().Interface()

				fglobalizer, ok := fref.(types.Contexter)

				if ok {
					fglobalizer.SetContext(dc.ctx)
				}

				ptrs[i] = fref
			} else {
				//				reflectValue := reflect.New(reflect.PtrTo(fval.Type()))

				//				fmt.Println(fval)
				//				reflectValue.Elem().Set(fval)

				//				ptrs[i] = reflectValue.Interface()
				//TODO: We cannot leave prts[i] nil
				logger.Errorf("The field %s is unaddressable", name)
			}

		} else {
			//TODO: We cannot leave prts[i] nil
			logger.Errorf("%s: column %s cannot be mapped to any field", reflect.TypeOf(m), cname)
		}
	}

	return ptrs
}

func (at *annotation) getFields(t reflect.Type) []*reflect.StructField {
	fields := make([]*reflect.StructField, 0, 10)

	for i := 0; i < t.NumField(); i++ {

		field := t.Field(i)
		if field.Anonymous {
			for _, f := range at.getFields(field.Type) {
				fields = append(fields, f)
			}
		} else {
			fields = append(fields, &field)
		}

		// if field.Anonymous == false {
		// 	fields = append(fields, &field)
		// }
	}

	return fields
}
