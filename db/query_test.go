package db_test

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/yaitoo/sparrow/db"
	"github.com/yaitoo/sparrow/db/model"

	"github.com/go-sql-driver/mysql"

	"testing"
)

var fakeConfig = model.Config{
	Database: model.Database{
		DSN:         "root:{passwd}@tcp(127.0.0.1:4010)/TransDB",
		Passwd:      "m_root_pwd",
		MaxConns:    10,
		MinConns:    10,
		MaxLifeTIme: 0 * time.Second,
	},
	Versions: []model.Version{
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

func TestQueryByKeySuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "172.31.0.11", Port: 3306, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	query.For("tran").Select("user_id", "amount").Where("id={id}", "id")

	vars := map[string]interface{}{"id": 1801928967192576}

	type Uid struct {
		UserID int64 //`sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	query.Vars(vars)
	uid := Uid{}
	err := query.Find(&uid)

	if err != nil {
		t.Error(err)
	}
	if uid.UserID != 456 || uid.Amount != 300 || uid.uuid != "" {
		t.Error("resutl not correct")
	}

}

func TestQueryByStructSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"id": 1801928967192576}

	type Uid struct {
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	uid := Uid{}
	query.For("tran").SelectModel(&uid).Where("id={id}", "id")
	query.Vars(vars)
	err := query.Find(&uid)
	if err != nil {
		t.Error(err)
	}
	if uid.Use_id != 456 || uid.Amount != 300 || uid.uuid != "" {
		t.Error("resutl not correct")
	}

}

func TestQueryByRawSqlSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"id": 1801928967192576}

	type Uid struct {
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	uid := Uid{}
	query.RawSQL("SELECT `tran`.user_id FROM tran").Where("id={id}", "id")
	query.Vars(vars)
	err := query.Find(&uid)
	if err != nil {
		t.Error(err)
	}
	if uid.Use_id != 456 || uid.Amount != 0 || uid.uuid != "" {
		t.Error("resutl not correct")
	}
}

func TestQueryRowByKeySuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"id": 1801928967192576}

	type Uid struct {
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	uid := Uid{}
	query.For("tran").SelectModel(&uid).Where("id={id}", "id")
	query.Vars(vars)

	err := query.QueryRow(func(row *sql.Row) error {
		if err := row.Scan(&uid.Use_id, &uid.Amount); err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	if uid.Use_id != 456 || uid.Amount != 300 || uid.uuid != "" {
		t.Error("resutl not correct")
	}
}

func TestQueryRowsByKeySuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	/* dsnMap := make(map[string]model.DSN)
	dsnMap["db00"] = model.DSN{NodeName: "db00", Host: "127.0.0.1", Port: 4006, Db: "TransDB", Username: "root", Password: "m_root_pwd"}
	*/
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"id": 1801928967192576}

	type Uid struct {
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	uid := Uid{}
	query.For("tran").SelectModel(&uid).Where("id={id}", "id")
	query.Vars(vars)

	err := query.Query(func(rows *sql.Rows) error {
		for rows.Next() {
			if err := rows.Scan(&uid.Use_id, &uid.Amount); err != nil {
				return err
			}

		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	if uid.Use_id != 456 || uid.Amount != 300 || uid.uuid != "" {
		t.Error("resutl not correct")
	}
}

func TestQuerySliceStruct(t *testing.T) {

	type Uid struct {
		ID   int64  `sql:"id"`
		Name string `sql:"name"`
		uuid string
	}

	uid := make([]Uid, 0)
	config := mysql.Config{
		User:                 "root",
		Passwd:               "m_root_pwd",
		Addr:                 "127.0.0.1:4006",
		Net:                  "tcp",
		DBName:               "TransDB",
		AllowNativePasswords: true,
	}
	dbi, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		t.Error(err)
	}
	defer dbi.Close()

	rows, err := dbi.Query("select * from test")
	err = db.ScanAll(rows, &uid, false)
	if err != nil {
		t.Error(err)
	}
	if len(uid) < 1 {
		t.Error(errors.New(""))
	}

}

func TestFindPrimptive(t *testing.T) {
	var familyID string

	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"scheduleid": "0208B1E1-4F72-4545-8722-E220894526BB"}

	query.For("cleanschedule").Select("FamilyId").Where("CleanScheduleId={scheduleid}", "scheduleid")
	query.Vars(vars)

	err := query.Find(&familyID)
	if err != nil {
		t.Error(err)
	}
}

func TestFindPrimptiveSlice(t *testing.T) {
	var cleanSchduleID []string

	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"familyid": "91b18f1f-4ef8-4066-97c4-28daea585db5"}

	query.For("cleanschedule").Select("CleanScheduleId").Where("FamilyId={familyid}", "familyid")
	query.Vars(vars)

	err := query.Find(&cleanSchduleID)
	if err != nil {
		t.Error(err)
	}
	if len(cleanSchduleID) == 0 {
		t.Fail()
	}
}

func TestFindCustomerStructure(t *testing.T) {
	type CleanSchedule struct {
		CleanScheduleID string
		CleanDateTime   time.Time
		FamilyID        string
		CleanItem       int
	}
	var familyID CleanSchedule

	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"scheduleid": "0208B1E1-4F72-4545-8722-E220894526BB"}

	query.For("cleanschedule").Select("FamilyId", "CleanScheduleId").Where("CleanScheduleId={scheduleid}", "scheduleid")
	query.Vars(vars)

	err := query.Find(&familyID)
	if err != nil {
		t.Error(err)
	}
}

func TestFindCustomerStructureSlice(t *testing.T) {
	type CleanSchedule struct {
		CleanScheduleID string // `sql:"CleanScheduleId"`
		CleanDateTime   time.Time
		FamilyID        string // `sql:"FamilyId"`
		CleanItem       int
	}
	var familyID []CleanSchedule = make([]CleanSchedule, 0)

	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"familyid": "91b18f1f-4ef8-4066-97c4-28daea585db5"}

	query.For("cleanschedule").Select("FamilyId", "CleanScheduleId").Where("FamilyId={familyid}", "familyid")
	query.Vars(vars)

	err := query.Find(&familyID)
	if err != nil {
		t.Error(err)
	}
	if len(familyID) == 0 {
		t.Fail()
	}
}

func TestQueryNoRecordSuccess(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())

	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)
	query := dbCtx.NewQuery(cancelCtx)
	vars := map[string]interface{}{"id": 1}

	type Uid struct {
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		uuid   string
	}
	uid := Uid{}
	query.For("test1").RawSQL("SELECT `test1`.user_id FROM test1").Where("id={id}", "id")
	query.Vars(vars)
	err := query.First(&uid)
	if err != sql.ErrNoRows {
		t.Error()
	}
	t.Log(err)

}

func TestQueryConcurrent(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 1000; i++ {
		go func() {
			query := dbCtx.NewQuery(cancelCtx)
			vars := map[string]interface{}{"id": 1801928967192576}

			type Uid struct {
				Use_id int64 `sql:"user_id"`
				Amount int   `sql:"amount"`
				// uuid   string
			}
			uid := Uid{}
			query.For("tran").SelectModel(&uid).Where("id={id}", "id")
			query.Vars(vars)
			err := query.QueryRow(func(row *sql.Row) error {
				if err := row.Scan(&uid.Use_id, &uid.Amount); err != nil {
					if err == sql.ErrNoRows {
						return nil
					}
					return err
				}
				return nil
			})
			if err != nil {
				t.Error(err)
			}
			defer wg.Done()
		}()
	}
	wg.Wait()
	t.Log("done")
}

func TestQueryIn(t *testing.T) {
	var (
		globalDb *db.Database
	)

	cancelCtx, _ := context.WithCancel(context.Background())
	globalDb = db.NewDatabase(cancelCtx, db.WithConfig(fakeConfig))
	dbCtx := globalDb.Open(cancelCtx)

	query := dbCtx.NewQuery(cancelCtx)
	//vars := map[string]interface{}{"id": []interface{}{2380033310064647, 2380033310064652}, "name": "test"}
	vars := map[string]interface{}{"id": 2380033310064647, "name": "test", "amount": []interface{}{300, 100}}

	type Uid struct {
		Id     int64 `sql:"id"`
		Use_id int64 `sql:"user_id"`
		Amount int   `sql:"amount"`
		// uuid   string
	}
	uids := []Uid{}
	query.For("tran").Select("id", "user_id").Where("id = {id}", "id").WhereOrIN("amount in ({amount})", "amount")
	query.Vars(vars)
	err := query.Find(&uids)
	if err != nil {
		t.Error(err)
	}
	t.Log("done")
}
