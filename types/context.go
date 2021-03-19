package types

import "context"

type appContext string

const myAppContext = appContext("ctx")

//Context  a request-scope context
type Context struct {
	RequestID int

	AcceptLanuage string

	TimeOffset int
	TimeJSON   string
	TimeLayout string
}

//WithContext 置入程序上下文对象
func WithContext(ctx context.Context, app Context) context.Context {
	return context.WithValue(ctx, myAppContext, app)
}

//FromContext 提取程序上下文
func FromContext(ctx context.Context) Context {
	val := ctx.Value(myAppContext)

	app, ok := val.(Context)
	if ok {
		return app
	}

	return Context{}
}
