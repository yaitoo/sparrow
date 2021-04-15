// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import (
	"fmt"
	"reflect"
	"sync"
)

//PrintfFunc  print function to print argument with type cast instead of reflect
type PrintfFunc func(source string, verb FormatVerb, flags FormatFlags, arg interface{}) string

var (
	printMutex sync.RWMutex
	printFuncs = make(map[string]PrintfFunc)
)

//RegisterPrintfFunc Register custom print format function for verb and type
func RegisterPrintfFunc(v FormatVerb, t reflect.Type, f PrintfFunc) {
	printMutex.Lock()
	defer printMutex.Unlock()

	key := fmt.Sprintf("%s-%v", v, t)

	_, ok := printFuncs[key]

	if !ok {
		printFuncs[key] = f
	}
}

func defaultPrintfFunc(format string, verb FormatVerb, flags FormatFlags, arg interface{}) string {
	return fmt.Sprintf(format, arg)
}

func getPrintfFunc(v FormatVerb, t reflect.Type) PrintfFunc {
	printMutex.RLock()
	defer printMutex.RUnlock()
	key := fmt.Sprintf("%s-%v", v, t)
	f, ok := printFuncs[key]
	if ok {
		return f
	}
	return defaultPrintfFunc
}
