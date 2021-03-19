package micro

import (
	"context"
	"fmt"
	"time"
)

// RequestContext a request-scope context
type RequestContext struct {
	RequestID   string
	RequsetTime time.Time

	AcceptLanguage string

	//types.Context
}

type ctxkey string

const requestContextKey = ctxkey("spx")

var ctx, cancel = context.WithCancel(context.Background())

//NewContext return root context of micro
func NewContext() context.Context {
	return ctx
}

//CancelContext cancel all contexts derived from micro context
func CancelContext() {
	cancel()
}

// WithContext returns a new Context with RequestConext.
func WithContext(ctx context.Context, rctx RequestContext) context.Context {
	return context.WithValue(ctx, requestContextKey, rctx)
}

// FromContext extracts the RequestContext from ctx, if present.
func FromContext(ctx context.Context) (RequestContext, bool) {
	rctx, ok := ctx.Value(requestContextKey).(RequestContext)
	return rctx, ok
}

// WithValues returns a copy of parent in which the value associated with values
func WithValues(parent context.Context, values map[string]string) context.Context {
	return &valuesCtx{parent, values}
}

// A valuesCtx carries a key-value pair. It implements Value for that key and
// delegates all other calls to the embedded Context.
type valuesCtx struct {
	context.Context
	values map[string]string
}

func (c *valuesCtx) String() string {
	return fmt.Sprintf("%v.WithValues(%#v)", c.Context, c.values)
}

func (c *valuesCtx) Value(key interface{}) interface{} {
	if c.values != nil {
		k, ok := key.(string)
		if ok {
			v, ok := c.values[k]
			if ok {
				return v
			}
		}

	}

	return c.Context.Value(key)
}
