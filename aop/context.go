package aop

import (
	"context"
	"time"
)

//Context aop上下文
type Context struct {
	ctx context.Context

	//	Config *Config

	currentHandlerIndex int

	FuncMetadata *FuncMetadata
	FuncInArgs   FuncInArgs

	Values map[interface{}]interface{}
}

//FuncInArgs 函数参数修正集合
type FuncInArgs []interface{}

//Deadline see context.Context.Deadline
func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.ctx.Deadline()
}

//Done see context.Context.Done
func (ctx *Context) Done() <-chan struct{} {
	return ctx.ctx.Done()
}

//Err see context.Context.Err
func (ctx *Context) Err() error {
	return ctx.ctx.Err()
}

//Value see context.Context.Value
func (ctx *Context) Value(key interface{}) interface{} {
	v, ok := ctx.Values[key]

	if ok {
		return v
	}

	return ctx.ctx.Value(key)
}
