// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package fmt

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

//PrintFunc  print function to print argument with type cast instead of reflect
type PrintFunc func(source string, verb FormatVerb, flags FormatFlags, arg interface{}) string

var (
	printMutex sync.RWMutex
	printFuncs = make(map[string]PrintFunc)
)

//RegisterPrintFunc Register custom print function for verb and type
func RegisterPrintFunc(ctx context.Context, v FormatVerb, t reflect.Type, f PrintFunc) {
	printMutex.Lock()
	defer printMutex.Unlock()

	key := fmt.Sprintf("%s-%v", v, t)

	_, ok := printFuncs[key]

	if !ok {
		printFuncs[key] = f
	}
}

func defaultPrintFunc(format string, verb FormatVerb, flags FormatFlags, arg interface{}) string {
	return fmt.Sprintf(format, arg)
}

func getPrintFunc(v FormatVerb, t reflect.Type) PrintFunc {
	printMutex.RLock()
	defer printMutex.RUnlock()
	key := fmt.Sprintf("%s-%v", v, t)
	f, ok := printFuncs[key]
	if ok {
		return f
	}
	return defaultPrintFunc
}
