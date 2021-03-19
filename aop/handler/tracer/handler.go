package tracer

import (
	"time"

	"github.com/rs/xid"
	"github.com/yaitoo/sparrow/aop"
)

//DefaultHandler 默認跟蹤器
var DefaultHandler = &Handler{callback: make([]func(ctx Context), 0, 2)}

//Context 跟蹤上下文
type Context struct {
	//ID  調用ID
	ID string
	//Func 函數名稱
	Func string
	//Before 進入時間
	Before time.Time
	//After  退出時間
	After time.Time

	ElapsedTime time.Duration
}

//Handler 函數調用跟蹤器
type Handler struct {
	callback []func(ctx Context)
}

func genTraceID() string {
	return xid.New().String()
}

var traceidKey = ctxkey("traceID")
var traceKey = ctxkey("traceCtx")

//Before 執行前
func (t *Handler) Before(ctx *aop.Context) error {
	v := ctx.Value(traceidKey)
	traceID := ""
	if v != nil {
		traceID, _ = v.(string)
	}

	if traceID == "" {
		traceID = genTraceID()
		ctx.Values[traceidKey] = traceID
	}

	tctx := &Context{
		ID:     traceID,
		Func:   ctx.FuncMetadata.Name,
		Before: time.Now(),
	}

	ctx.Values["tctx_"+ctx.FuncMetadata.Name] = tctx

	return nil

}

//After 執行後
func (t *Handler) After(ctx *aop.Context) {

	v, ok := ctx.Values["tctx_"+ctx.FuncMetadata.Name]

	if ok {
		tctx, ok := v.(*Context)

		if ok {
			tctx.After = time.Now()
			tctx.ElapsedTime = tctx.After.Sub(tctx.Before)

			for _, callback := range t.callback {
				go callback(*tctx)
			}

		}
	}
}

//OnInvoked 函數調用跟蹤結束
func (t *Handler) OnInvoked(callback func(ctx Context)) {
	if t == nil {
		return
	}

	if t.callback == nil {
		t.callback = make([]func(ctx Context), 0, 2)
	}

	if callback != nil {
		t.callback = append(t.callback, callback)
	}

}
