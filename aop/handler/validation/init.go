package validation

import "github.com/yaitoo/sparrow/log"

var (
	logger = log.NewLogger("authn")
	//DefaultHandler 默認AOP攔截器
	DefaultHandler *Handler
)

func init() {
	DefaultHandler = NewHandler()
}
