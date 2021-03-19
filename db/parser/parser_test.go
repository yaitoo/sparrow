package parser_test

import (
	"testing"

	"github.com/yaitoo/sparrow/db/model"
	"github.com/yaitoo/sparrow/db/parser"
	"github.com/yaitoo/sparrow/db/shardingId"
)

var commandfakeConfig = model.Config{
	Versions: []model.Version{
		{
			Version: 0,
			Rules: model.Rule{
				Databases: []model.Database{
					model.Database{
						DSN:    "root:{passwd}@tcp(127.0.0.1:4006)/TransDB",
						Passwd: "m_root_pwd",
					},
				},
				Tables: []model.Table{},
			},
		},
		{
			Version: 1,
			Rules: model.Rule{
				Databases: []model.Database{
					model.Database{
						DSN:    "root:{passwd}@tcp(127.0.0.1:4006)/TransDB",
						Passwd: "m_root_pwd",
					},
				},
				Tables: []model.Table{
					model.Table{
						Name:     "tran",
						Key:      "id",
						TimeZone: "",
						Tags: []model.Tag{
							model.Tag{
								ID: 0,
								Names: []string{
									"deposit"},
								Amount: 5,
								Date:   "month",
							},
						},
					},
				},
			},
		},
	},
}

func TestParseWithVarMapSuccess(t *testing.T) {
	sqlStr := "INSERT INTO tran (user_id, id)  VALUES ({user_id},{id});"
	var id int64 = 1154877550153957376
	var idObj = shardingId.ParseId(id)

	var expected parser.SqlStringAndEnitityData = parser.SqlStringAndEnitityData{
		EntityData: make(parser.EntityKeyPair),
		SqlString:  "insert into tran_00_2019_11_00(user_id, id) values (?, ?)",
	}
	expected.EntityData["tran"] = parser.AliasEntityKeyPair{
		EntityName: "tran",
		Alias:      "",
		Key: parser.KeyAndValue{
			Key:   "id",
			Value: idObj,
		},
	}

	vars := map[string]interface{}{"id": id, "user_id": 456}
	cmds, err := parser.ParseWithVarMap(commandfakeConfig, sqlStr, vars)
	if err != nil {
		t.Error(err)
	}

	result := cmds
	if result.SqlString != expected.SqlString {
		t.Error("not equal")
	}
}
