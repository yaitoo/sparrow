package model

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/db/util"
)

var secretKey string
var passwdTemplate = "{passwd}"

const defaultConns int = 20
const defaultMaxLifetime time.Duration = 0

type Database struct {
	DSN         string        `yaml:"dsn"`
	Passwd      string        `yaml:"passwd"`
	MaxConns    int           `yaml:"max_conns"`
	MinConns    int           `yaml:"min_conns"`
	MaxLifeTIme time.Duration `yaml:"max_lifetime"`
}

func (d *Database) Validate() bool {
	if d.IsNullObject() == true {
		return false
	}
	return true
}

func (d *Database) IsNullObject() bool {
	return d.DSN == ""
}

func (n *Database) String() string {
	result, _ := json.Marshal(n)
	return string(result)
}

func (n Database) ConnStr() (string, error) {
	var result = n.DSN
	if secretKey != "" {
		plaintPwd, err := util.AesDecrypt(n.Passwd, secretKey)
		if err != nil {
			return "", err
		}

		return strings.Replace(result, passwdTemplate, plaintPwd, 1), nil
	}

	return strings.Replace(result, passwdTemplate, n.Passwd, 1), nil
}

func (d *Database) SetMaxConns(connAmout int) {
	if connAmout == 0 {
		d.MaxConns = defaultConns
		return
	}
	d.MaxConns = connAmout
}
func (d *Database) SetMinConns(connAmout int) {
	if connAmout == 0 {
		d.MinConns = defaultConns
		return
	}
	d.MinConns = connAmout
}
func (d *Database) SetMaxLifeTime(duration time.Duration) {
	if duration == 0 {
		d.MaxLifeTIme = defaultMaxLifetime
		return
	}
	d.MaxLifeTIme = duration
}
