package authz

import "context"

type ctxkey string

var roleKey = ctxkey("$role")

//NewContext 创建验证上下文
func NewContext(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, roleKey, role)
}

//FromContext 从context中提取登录用户
func FromContext(ctx context.Context) string {
	v := ctx.Value(roleKey)

	if v == nil {
		return ""
	}

	id, ok := v.(string)
	if ok {
		return id
	}

	return id
}
