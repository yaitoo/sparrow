package db

import (
	"github.com/yaitoo/sparrow/db/model"
)

//Option 变数设定
type Option func(d *Database)

//WithConfig 指定设定档
func WithConfig(c model.Config) Option {
	return func(d *Database) {
		d.config = func() model.Config {
			return c
		}
	}
}
