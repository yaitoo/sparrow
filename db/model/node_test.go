package model_test

import (
	"testing"

	"github.com/yaitoo/sparrow/db/model"
)

var version string

func TestNode(t *testing.T) {
	t.Log(version)
	node := &model.Database{
		DSN:    "root:{passwd}@tcp(127.0.0.1:4006)/TransDB",
		Passwd: "m_root_pwd",
	}
	conn, err := node.ConnStr()
	if err != nil {
		t.Fatal(err)
	}
	if conn != "root:m_root_pwd@tcp(127.0.0.1:4006)/TransDB" {
		t.Fatal()
	}
}
