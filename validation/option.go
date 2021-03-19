package validation

import (
	"github.com/BurntSushi/toml"
	"github.com/yaitoo/sparrow/validation/models"
)

//Option 变数设定
type Option func(ctx *Context)

//WithForms 指定表单设定
func WithForms(raw string) Option {
	return func(ctx *Context) {
		if len(raw) > 0 {
			forms := make(map[string]models.Form)
			if _, err := toml.Decode(raw, &forms); err != nil {
				logger.Warnln(err)
			}

			ctx.forms = forms
		}
	}
}
