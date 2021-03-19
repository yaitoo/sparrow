package authn

import "context"

type ctxkey string

var identityUserKey = ctxkey("$authn_id")

//NewContext 创建验证上下文
func NewContext(ctx context.Context, user IdentityUser) context.Context {
	return context.WithValue(ctx, identityUserKey, &user)
}

//FromContext 从context中提取登录用户
func FromContext(ctx context.Context) *IdentityUser {
	v := ctx.Value(identityUserKey)

	if v == nil {
		return nil
	}

	id, ok := v.(*IdentityUser)
	if ok {
		return id
	}

	return id
}
