package db_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yaitoo/sparrow/db"
	"github.com/yaitoo/sparrow/db/model"
	"github.com/yaitoo/sparrow/db/shardingId"
)

var commandfakeConfig = model.Config{
	Database: model.Database{
		DSN:         "root:{passwd}@tcp(127.0.0.1:4010)/TransDB",
		Passwd:      "m_root_pwd",
		MaxConns:    10,
		MinConns:    5,
		MaxLifeTIme: 0, //1 * time.Second,
	}, Versions: []model.Version{
		{
			Version: 0,
			Rules: model.Rule{
				Databases: []model.Database{
					model.Database{
						DSN:         "root:{passwd}@tcp(127.0.0.1:4010)/TransDB",
						Passwd:      "m_root_pwd",
						MaxConns:    5,
						MinConns:    5,
						MaxLifeTIme: 0,
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
								ID:     0,
								Names:  []string{"cp"},
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

func TestDeleteSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)
	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	com := dbCtx.NewCommand(cancelCtx)

	vars := map[string]interface{}{"id": 1801928967192576}

	cmd := com.For("tran").Where("id = {id}", "id").Delete()
	cmd.Vars(vars)
	_, err := cmd.Exec()
	if err != nil {
		t.Error(err)
	}
}

func TestInsertSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)
	shardingId.TimeNow = time.Date(2019, time.October, 20, 10, 15, 30, 0, time.UTC)
	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	com := dbCtx.NewCommand(cancelCtx)
	// id, _ := globalDb.NewID(cancelCtx, "tran", "deposit")

	_, err := com.For("tran").Insert("id", 1801928967192576).Insert("user_id", 456).Exec()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	com := dbCtx.NewCommand(cancelCtx)

	vars := map[string]interface{}{"id": 1801928967192576}

	cmd := com.For("tran").Update("amount", 300).Where("id = {id}", "id")
	cmd.Vars(vars)
	_, err := cmd.Exec()
	if err != nil {
		t.Error(err)
	}
}

func TestTxCrossDB(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	com := dbCtx.NewCommand(cancelCtx)

	vars := map[string]interface{}{"id": 1629860229283840}
	com.For("tran").Update("amount", 300).Where("id = {id}", "id")
	com.Vars(vars)
	/* 	cmd.Vars(vars) */
	err := dbCtx.Begin()
	_, err = com.Exec()
	if err != nil {
		t.Error(err)
	}

	_, err = com.For("tran").Insert("id", 1629860230332417).Insert("user_id", 456).Exec()
	if err != nil {
		dbCtx.Rollback()
		if err.Error() == "cross database" {
			t.Log(err)
		} else {
			t.Error(err)
		}
	} else {
		dbCtx.Commit()
	}
}

func TestTxDB(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	com := dbCtx.NewCommand(cancelCtx)

	vars := map[string]interface{}{"id": 1629860230332417}

	cmd := com.For("order").Where("id = {id}", "id").Delete()
	cmd.Vars(vars)
	_, err := cmd.Exec()

	dbCtx = globalDb.Open(cancelCtx)
	com = dbCtx.NewCommand(cancelCtx)

	vars = map[string]interface{}{"id": 1629860229283840}
	com.For("order").Insert("id", 1629860230332417).Insert("user_id", 456)

	/* 	cmd.Vars(vars) */
	err = dbCtx.Begin()
	_, err = com.Exec()
	if err != nil {
		t.Error(err)
	}

	com.For("tran").Update("amount", 300).Where("id = {id}", "id")
	com.Vars(vars)
	if err != nil {
		dbCtx.Rollback()
		if err.Error() == "cross database" {
			t.Log(err)
		} else {
			t.Error(err)
		}
	} else {
		dbCtx.Commit()
	}
}

func TestMultiTxDB(t *testing.T) {
	var (
		globalDb *db.Database
		err      error
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	gwg := sync.WaitGroup{}
	gwg.Add(1)
	for k := 0; k < 1; k++ {
		go func() {
			defer gwg.Done()
			dbCtx := globalDb.Open(cancelCtx)

			wg := sync.WaitGroup{}
			wg.Add(100)

			for i := 0; i < 100; i++ {
				go func(idx int, odb *db.Database) {

					oid, _ := odb.NewID(cancelCtx, "order", "cp")
					com := dbCtx.NewCommand(cancelCtx)
					com.For("order").Insert("id", oid).Insert("user_id", 333)
					err = dbCtx.Begin()
					_, err = com.Exec()

					if err != nil {
						fmt.Println(err)
						// dbCtx.Rollback()

					} else {
						time.Sleep(1 * time.Millisecond)
						dbCtx.Commit()
					}
					defer wg.Done()
				}(i, globalDb)
			}
			wg.Wait()
		}()
	}
	gwg.Wait()
	globalDb.Close()
	if err != nil {
		t.Fail()
	}
	t.Log()
}

func TestRaw(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(commandfakeConfig))
	dbCtx := globalDb.Open(cancelCtx)

	updCol := map[string]interface{}{
		"qrcode": 1,
		"status": 2,
	}

	vars := map[string]interface{}{
		"member_id":  2,
		"app_dtl_id": 3,
	}

	cmd := dbCtx.NewCommand(cancelCtx).
		For("payment_member").
		Updates(updCol).
		Where("member_id = {member_id}", "member_id").
		Where("app_dtl_id = {app_dtl_id}", "app_dtl_id").
		Vars(vars)

	_, err := cmd.Exec()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	t.Log()
}
