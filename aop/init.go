package aop

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/yaitoo/sparrow/log"
)

var (
	cfg    *Config
	logger = log.NewLogger("aop")
)

func init() {
	cfg = &Config{}
	filePath, err := filepath.Abs("./conf.d/aop.toml")
	if err != nil {
		//panic(err)
		logger.Warnln(err)
	}

	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		logger.Errorln(err)
	}
}
