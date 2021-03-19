package aop

import (
	"reflect"
	"sync"
)

var (
	valuers      map[reflect.Type]CustomValueHandle = make(map[reflect.Type]CustomValueHandle)
	valuersMutex sync.RWMutex
)

//Valuer 将各种型别的变数转化成健值对
type Valuer interface {
	Value() map[string]interface{}
}

//CustomValueHandle 自定义Valuer展开函数
type CustomValueHandle func(obj interface{}) map[string]interface{}

//RegisterCustomValueHandle 注册自定义Valuer展开函数
func RegisterCustomValueHandle(t reflect.Type, valuer CustomValueHandle) {
	valuersMutex.Lock()
	defer valuersMutex.Unlock()

	valuers[t] = valuer
}

//getCustomValueHandle 获取指定型别的Valuer自定义展开函数
func getCustomValueHandle(t reflect.Type) CustomValueHandle {
	valuersMutex.RLock()
	defer valuersMutex.RUnlock()

	v, ok := valuers[t]

	if ok {
		return v
	}

	return nil
}

//IndirectType 获取真正型别
func IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	return reflectType
}
