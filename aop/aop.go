package aop

import (
	"context"
	"reflect"
	"runtime"
)

var handlers = make([]Handler, 0, 10)

//Use 註冊中間件
func Use(middlewares ...Handler) {
	for _, m := range middlewares {
		handlers = append(handlers, m)
	}

}

//RegisterNamesIn 手工注册函数的参数名称列表
func RegisterNamesIn(fn interface{}, names ...string) {

	fnName := cfg.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())

	cfg.setNamesIn(fnName, names...)
}

//NewContext 創建AOP上下文
func NewContext(ctx context.Context) *Context {
	return &Context{ctx: ctx}
}

//Before 函式執行前
func Before(ctx context.Context, fn interface{}, args ...interface{}) (context.Context, error) {

	actx := NewContext(ctx)
	actx.FuncMetadata, actx.FuncInArgs = loadFuncMetadata(fn, args)
	actx.Values = make(map[interface{}]interface{})

	//actx.FuncInArgs = func() []interface{} { return args }

	for i, handler := range handlers {
		actx.currentHandlerIndex = i
		if err := handler.Before(actx); err != nil {
			return actx, err
		}
	}

	return actx, nil
}

//After 函式執行後
func After(ctx context.Context) {

	max := len(handlers)
	if max > 0 {

		actx, ok := ctx.(*Context)
		if ok {
			for i := actx.currentHandlerIndex; i >= 0; i-- {
				handlers[i].After(actx)
			}
		}
	}
}
