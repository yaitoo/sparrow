package authn

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/yaitoo/sparrow/log"
)

var (
	cfg    *Config
	logger = log.NewLogger("authn")
	//DefaultHandler 默認AOP攔截器
	DefaultHandler *Handler
)

func init() {
	cfg = &Config{}
	filePath, err := filepath.Abs("./conf.d/authn.toml")
	if err != nil {
		//panic(err)
		logger.Warnln(err)
		return
	}

	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		logger.Errorln(err)
		return
	}

	DefaultHandler = NewHandler(cfg)

}
