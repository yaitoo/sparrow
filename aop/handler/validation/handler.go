package validation

import (
	"github.com/yaitoo/sparrow/aop"
	"github.com/yaitoo/sparrow/micro"
	"github.com/yaitoo/sparrow/validation"
)

//Handler 函數調用跟蹤器
type Handler struct {
}

//NewHandler 创建Handler
func NewHandler() *Handler {
	h := &Handler{}

	return h
}

//Before 執行前
func (h *Handler) Before(ctx *aop.Context) error {
	if h == nil {
		return nil
	}

	formName := ctx.FuncMetadata.Name
	fieldNames := ctx.FuncMetadata.FixedNamesIn

	if len(fieldNames) > 0 {
		vg := validation.NewContext(ctx)
		vg.ValidateGroup(formName).
			WithNames(fieldNames...).
			WithValues(ctx.FuncInArgs...)

		if vg.IsValid() == false {
			return micro.Throw(ctx, micro.ErrBadRequest, vg.Error())
		}
	}

	return nil

}

//After 執行後
func (h *Handler) After(ctx *aop.Context) {
}
