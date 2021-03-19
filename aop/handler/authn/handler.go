package authn

import (
	"github.com/yaitoo/sparrow/aop"
	"github.com/yaitoo/sparrow/micro"
)

//Handler 函數調用跟蹤器
type Handler struct {
	cfg *Config
}

//NewHandler 创建Handler
func NewHandler(c *Config) *Handler {
	h := &Handler{}
	h.cfg = cfg

	return h
}

//Before 執行前
func (h *Handler) Before(ctx *aop.Context) error {
	if h == nil || h.cfg == nil {
		return nil
	}

	if !h.cfg.IsIdentityUser(FromContext(ctx), ctx.FuncMetadata.Name) {
		return micro.Throw(ctx, micro.ErrUnauthorized, "请先登录")
	}

	return nil

}

//After 執行後
func (h *Handler) After(ctx *aop.Context) {
}
