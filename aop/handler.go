package aop

//Handler 請求攔截器
type Handler interface {
	//Before 執行前
	Before(ctx *Context) error
	//After 執行後
	After(ctx *Context)
}
