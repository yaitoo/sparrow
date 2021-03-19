package aop

import (
	"strings"
	"sync"
)

//Config AOP設定
type Config struct {
	sync.RWMutex
	//Pkg 包名替換規則
	Pkg map[string]string
	//NamesIn 函数参数名称
	NamesIn map[string][]string
}

//Replace 替換包名
func (cfg *Config) Replace(name string) string {
	if cfg == nil || cfg.Pkg == nil {
		return name
	}

	for k, v := range cfg.Pkg {
		if strings.HasPrefix(name, k) {
			return v + string(name[len(k):])
		}
	}

	return name
}

//getNamesIn 读取指定函数的参数名称列表
func (cfg *Config) getNamesIn(fnName string) ([]string, bool) {
	cfg.RLock()
	defer cfg.RUnlock()

	v, ok := cfg.NamesIn[fnName]

	return v, ok
}

//setNamesIn 设定指定函数的参数名称列表
func (cfg *Config) setNamesIn(fnName string, names ...string) {
	cfg.Lock()
	defer cfg.Unlock()

	cfg.NamesIn[fnName] = names
}
