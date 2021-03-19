package validation

import (
	"github.com/BurntSushi/toml"
	"github.com/yaitoo/sparrow/config"
	"github.com/yaitoo/sparrow/log"
	"github.com/yaitoo/sparrow/validation/models"
	"github.com/yaitoo/sparrow/validation/validator"
)

const (
	defaultLang       = "zh_CN"
	defaultLangCtxKey = "lang"
)

var (
	logger = log.NewLogger("validation")
	//forms 表单验证集合
	forms = make(map[string]models.Form)
)

func init() {
	file, raw, err := config.Open("./conf.d/validation.toml")
	if err != nil {
		logger.Warnln(err)
	}
	//继续添加监听事件，补足文件即可生效
	file.OnFileChanged(loadForms)

	loadForms(raw)

	registerValidators()
}

func loadForms(raw []byte) {
	if len(raw) > 0 {
		if _, err := toml.Decode(string(raw), &forms); err != nil {
			logger.Warnln(err)
		}
	}
}

func registerValidators() {
	Register("Required", &validator.Required{})
	Register("Min", &validator.Min{})
	Register("Max", &validator.Max{})
	Register("Range", &validator.Range{})
	Register("MinSize", &validator.MinSize{})
	Register("MaxSize", &validator.MaxSize{})
	Register("Length", &validator.Length{})
	Register("Alpha", &validator.Alpha{})
	Register("Numeric", &validator.Numeric{})
	Register("AlphaNumeric", &validator.AlphaNumeric{})
	Register("Match", &validator.Match{})
	Register("NoMatch", &validator.NoMatch{})
	Register("AlphaDash", &validator.AlphaDash{})
	Register("Email", &validator.Email{})
	Register("Base64", &validator.Base64{})
	Register("IP", &validator.IPv4{})
	Register("Mobile", &validator.Mobile{})
	Register("Tel", &validator.TEL{})
	Register("Phone", &validator.Phone{})
	Register("ZipCode", &validator.ZipCode{})
}
