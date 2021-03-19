package aop

import (
	"reflect"
	"runtime"
	"strings"
	"sync"
)

var (
	funcMetadata = sync.Map{}
)

//FuncMetadata func反射後得到的元數據
type FuncMetadata struct {
	//Type 型別信息
	Type reflect.Type
	//Name 完整名稱包含包地址
	Name string
	//NumIn 變數數量
	NumIn int
	//NamesIn 變數名稱
	NamesIn []string
	//FixedNamesIn -/+ 修正后的参数列表
	FixedNamesIn []string
}

func loadFuncMetadata(fn interface{}, args []interface{}) (*FuncMetadata, []interface{}) {
	if fn == nil {
		return nil, nil
	}

	t := reflect.TypeOf(fn)

	// i, ok := funcMetadata.Load(t)
	// if ok {
	// 	m, ok := i.(*FuncMetadata)
	// 	if ok {
	// 		return m
	// 	}
	// }

	m := &FuncMetadata{}

	m.Type = t
	m.NumIn = len(args)
	m.Name = cfg.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())

	namesIn, ok := cfg.getNamesIn(m.Name)

	if !ok {
		namesIn = make([]string, m.NumIn)
		for i := 0; i < m.NumIn; i++ {
			namesIn[i] = "+"
		}
	}

	m.NamesIn = namesIn

	//Fixed -/+
	fixedNamesIn := make([]string, 0, m.NumIn)

	//TODO 改成缓存读取函数代替反射获取，提高效能
	funcInArgs := make([]interface{}, 0, m.NumIn)

	for i, name := range namesIn {
		if name == "-" { //忽略

		} else if name == "+" { //展开map/struct
			arg := args[i]

			valuer, ok := arg.(Valuer)
			//本身实现了aop.Valuer
			if ok {
				values := valuer.Value()
				for k, v := range values {
					fixedNamesIn = append(fixedNamesIn, strings.ToLower(k))
					funcInArgs = append(funcInArgs, v)
				}
				continue
			}

			// v := reflect.ValueOf(arg)

			// if v.Kind() == reflect.Ptr {
			// 	v = v.Elem()
			// }

			v := reflect.Indirect(reflect.ValueOf(arg))

			customValuer := getCustomValueHandle(t)

			//找到自定义的Valuer
			if customValuer != nil {
				values := customValuer(arg)
				for k, v := range values {
					fixedNamesIn = append(fixedNamesIn, strings.ToLower(k))
					funcInArgs = append(funcInArgs, v)
				}
				continue
			}

			//泛型模型
			switch v.Kind() {
			case reflect.Struct:
				//	for i :=0; i< t.FiledNum
				n := v.NumField()
				t := v.Type()
				for i := 0; i < n; i++ {

					fieldValue := v.Field(i)
					if fieldValue.CanInterface() {
						field := t.Field(i)

						tag := field.Tag.Get("val")
						if tag == "-" {
							continue
						} else if tag != "" {
							fixedNamesIn = append(fixedNamesIn, strings.ToLower(tag))
							funcInArgs = append(funcInArgs, fieldValue.Interface())
						} else {
							fixedNamesIn = append(fixedNamesIn, strings.ToLower(field.Name))
							funcInArgs = append(funcInArgs, fieldValue.Interface())
						}
					}

				}

			case reflect.Map:
				val := reflect.ValueOf(arg)
				for _, key := range val.MapKeys() {

					v := val.MapIndex(key).Interface()

					switch v.(type) {
					// case int:
					// 	fmt.Println(e, t)
					case string:
						//只支持string的值
						fixedNamesIn = append(fixedNamesIn, strings.ToLower(key.String()))
						funcInArgs = append(funcInArgs, v)

						// default:
						// 	fmt.Println("not found")

					}
				}
			}

		} else {
			fixedNamesIn = append(fixedNamesIn, strings.ToLower(name))
			funcInArgs = append(funcInArgs, args[i])
		}
	}

	m.FixedNamesIn = fixedNamesIn

	//fmt.Println(m.Name, m.NamesIn, m.FixedNamesIn)

	//funcMetadata.Store(t, m)

	return m, funcInArgs
}
