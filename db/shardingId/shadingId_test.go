package shardingId_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yaitoo/sparrow/db/model"
	"github.com/yaitoo/sparrow/db/shardingId"
)

var commandfakeConfig = &model.Config{
	Database: model.Database{
		DSN:    "root:{passwd}@tcp(127.0.0.1:4006)/TransDB",
		Passwd: "m_root_pwd",
	},
	Versions: []model.Version{
		{
			Version: 0,
			Rules: model.Rule{
				Databases: []model.Database{
					model.Database{
						DSN:    "root:{passwd}@tcp(127.0.0.1:4006)/TransDB",
						Passwd: "m_root_pwd",
					},
					model.Database{
						DSN:    "root:{passwd}@tcp(127.0.0.2:4006)/TransDB",
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
									"deposit",
								},
								Amount: 5,
								Date:   "month",
							},
						},
					},
					model.Table{
						Name:     "order",
						Key:      "id",
						TimeZone: "",
						Tags: []model.Tag{
							model.Tag{
								ID: 0,
								Names: []string{
									"cp",
								},
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

func TestSharding(t *testing.T) {
	targetTIme := time.Date(2019, 10, 7, 13, 10, 5, 0, time.UTC)
	shardingId.TimeNow = targetTIme
	gen1, err := shardingId.NewIdGenerator("tran", "deposit", commandfakeConfig.GetNewestVersion(), 0)

	if err != nil {
		t.Fatal(err)
	}

	id0, err := gen1.NextID()
	if err != nil {
		t.Fatal(err)
	}
	id1, err := gen1.NextID()
	if err != nil {
		t.Fatal(err)
	}
	id0Object := shardingId.ParseId(id0)
	idObject := shardingId.ParseId(id1)
	if idObject.Time.UTC() != targetTIme.UTC() {
		t.Fatal()
	}
	fmt.Println(id0Object.Time.UTC())
	fmt.Println(idObject.Time.UTC())
	t.Log(id1)
}

func TestSubSharding(t *testing.T) {
	targetTIme := time.Date(2019, 10, 7, 13, 10, 5, 0, time.UTC)
	shardingId.TimeNow = targetTIme
	gen1, err := shardingId.NewIdGenerator("tran", "deposit", commandfakeConfig.GetNewestVersion(), 0)

	if err != nil {
		t.Fatal(err)
	}

	id1, err := gen1.NextID()
	if err != nil {
		t.Fatal(err)
	}

	idObject := shardingId.ParseId(id1)
	fmt.Println(idObject)
	gen2, err := shardingId.NewSubId("order", "cp", commandfakeConfig.GetNewestVersion(), id1)
	if err != nil {
		t.Fatal(err)
	}
	subid1, err := gen2.NextID()
	if err != nil {
		t.Fatal(err)
	}

	subid1Object := shardingId.ParseId(subid1)

	fmt.Println(subid1Object)

	t.Log(id1)

}

func TestShardingIdNotConfigured(t *testing.T) {
	targetTIme := time.Date(2019, 10, 7, 13, 10, 5, 0, time.UTC)
	shardingId.TimeNow = targetTIme
	gen1, err := shardingId.NewIdGenerator("tran12", "deposit12", commandfakeConfig.GetNewestVersion(), 0)

	if err != nil {
		t.Fatal(err)
	}

	id0, err := gen1.NextID()
	if err != nil {
		t.Fatal(err)
	}
	id0Object := shardingId.ParseId(id0)
	if id0Object.Time.UTC() != targetTIme.UTC() {
		t.Fatal()
	}
	fmt.Println(id0Object.Time.UTC())

	for i := 0; i < 10; i++ {
		id, err := gen1.NextID()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(id)
	}
}
