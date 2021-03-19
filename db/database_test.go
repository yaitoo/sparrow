package db_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/yaitoo/sparrow/db"
	"github.com/yaitoo/sparrow/db/model"
	"github.com/yaitoo/sparrow/db/shardingId"
)

const (
	// cEpoch : Monday, May 20, 2019 12:00:00 AM
	cEpoch = 1558310400

	// cBusinessNumberBIts : business number
	cBusinessNumberBIts = 3

	// cDbIdBIts : length of database
	cDbIdBIts = 4

	// cSequenceBits : Length of sequence
	cSequenceBits = 20

	// cTimeStampBits : Length of time
	cTimeStampBits = 30

	// cInstanceBits : length of instace
	cInstanceBits = 3

	// cAlgVerBits : length of alg version
	cAlgVerBits = 3

	cBusinessNumberShift = cDbIdBIts + cSequenceBits
	cTimeStampShift      = cBusinessNumberShift + cBusinessNumberBIts
	cInstanceShift       = cTimeStampShift + cTimeStampBits
	cAlgVersionShift     = cInstanceShift + cInstanceBits

	cSequenceMask       = -1 ^ (-1 << cSequenceBits) //0xfffff // 1048575
	cMaxDBusinessNumber = -1 ^ (-1 << cBusinessNumberBIts)
	cMaxDb              = -1 ^ (-1 << cDbIdBIts)
	cMaxTimestamp       = -1 ^ (-1 << cTimeStampBits)
	cMaxInstance        = -1 ^ (-1 << cInstanceBits)
	cMaxAlgVer          = -1 ^ (-1 << cAlgVerBits)
)

var databasefakeConfig = model.Config{
	Versions: []model.Version{
		{
			Version: 0,
			Rules: model.Rule{
				Databases: []model.Database{
					model.Database{
						DSN:    "root:{passwd}@tcp(127.0.0.1:4010)/TransDB",
						Passwd: "m_root_pwd",
					},
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
								ID:     0,
								Names:  []string{"deposit"},
								Amount: 5,
								Date:   "month",
							},
							model.Tag{
								ID:     1,
								Names:  []string{"withdtraw"},
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
							model.Tag{
								ID:     1,
								Names:  []string{"mg", "ag", "sb"},
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

/* func TestDatabaseInit(t *testing.T) {
	cancelCtx, _ := context.WithCancel(context.Background())
	database := db.NewDatabase(cancelCtx)
	if database == nil {
		t.Error()
	}
} */

/* func TestDatabaseSetConf(t *testing.T) {
	cancelCtx, _ := context.WithCancel(context.Background())
	database := db.NewDatabase(cancelCtx).SetConf(func() model.Versions { return databasefakeConfig })
	if database ==  {
		t.Error()
	}
} */

func TestDatabaseNewID(t *testing.T) {
	// shardingId.TimeNow = time.Now()
	cancelCtx, _ := context.WithCancel(context.Background())
	db.InstanceID = 0
	database := db.NewDatabase(cancelCtx, db.WithConfig(databasefakeConfig))

	pid1, err := database.NewID(cancelCtx, "tran", "deposit")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pid1)
	time.Sleep(2 * time.Second)
	pid2, err := database.NewID(cancelCtx, "tran", "deposit")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pid2)
	pid3, err := database.NewID(cancelCtx, "tran", "deposit")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(pid3)

	pid1o := shardingId.ParseId(pid1)
	pid2o := shardingId.ParseId(pid2)
	pid3o := shardingId.ParseId(pid3)
	fmt.Println(pid1o, pid2o, pid3o)
}

func TestGenNewID(t *testing.T) {
	// shardingId.TimeNow = time.Now()
	cancelCtx, _ := context.WithCancel(context.Background())
	db.InstanceID = 0
	database := db.NewDatabase(cancelCtx, db.WithConfig(databasefakeConfig))

	ticker := time.NewTicker(5 * time.Second)
	var idMap sync.Map
	for {
		select {
		case <-ticker.C:
			return
		default:
			go func() {
				pid0, _ := database.NewID(cancelCtx, "tran", "deposit")
				pid1, _ := database.NewID(cancelCtx, "tran", "deposit")
				pid2, _ := database.NewSubID(cancelCtx, "tran", "deposit", pid1)
				_, exist0 := idMap.Load(pid0)
				_, exist1 := idMap.Load(pid1)
				_, exist2 := idMap.Load(pid2)
				if exist0 || exist1 || exist2 {
					t.Failed()
					return
				}
				idMap.Store(pid0, struct{}{})
				idMap.Store(pid1, struct{}{})
				idMap.Store(pid2, struct{}{})
			}()
		}
	}
}

type TestAtomic struct {
	algVerId       int64
	instance       int64
	businessNumber int64
	dbId           int64
	lastTimeStamp  int64
	sequence       *atomic.Int64

	lock      *sync.Mutex
	dbCounter *atomic.Int32
	nodesCnt  int
	isSub     bool
}

type TestLock struct {
	algVerId       int64
	instance       int64
	businessNumber int64
	dbId           int64
	lastTimeStamp  int64
	sequence       int64

	lock      *sync.Mutex
	dbCounter int
	nodesCnt  int
	isSub     bool
}

func (iw *TestAtomic) NextID() (uID int64, err error) {
	if iw.isSub == false {
		iw.dbId = int64((iw.dbCounter.Load()) % int32(iw.nodesCnt))
		if iw.dbCounter.Load() == math.MaxInt32 {
			iw.dbCounter.Store(0)
		}
		iw.dbCounter.Inc()
	}

	ts := timeGen()
	if iw.sequence.Load() >= math.MaxInt64 {
		ts = timeReGen(ts)
		iw.sequence.Store(0)
	} else {
		iw.sequence.Inc()
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("ID Error")
		return 0, err
	}

	iw.lastTimeStamp = ts
	uID = iw.algVerId<<cAlgVersionShift | iw.instance<<cInstanceShift | (ts-cEpoch)<<cTimeStampShift | iw.businessNumber<<cBusinessNumberShift | iw.dbId<<cSequenceBits | iw.sequence.Load()
	return uID, nil
}

func (iw *TestLock) NextID() (uID int64, err error) {
	iw.lock.Lock()
	defer iw.lock.Unlock()

	if iw.isSub == false {
		iw.dbId = int64(int32(iw.dbCounter) % int32(iw.nodesCnt))
		if iw.dbCounter == math.MaxInt32 {
			iw.dbCounter = 0
		}
		iw.dbCounter++
	}

	ts := timeGen()
	if iw.sequence >= math.MaxInt64 {
		ts = timeReGen(ts)
		iw.sequence = 0
	} else {
		iw.sequence++
	}

	if ts < iw.lastTimeStamp {
		err = errors.New("ID Error")
		return 0, err
	}

	iw.lastTimeStamp = ts
	uID = iw.algVerId<<cAlgVersionShift | iw.instance<<cInstanceShift | (ts-cEpoch)<<cTimeStampShift | iw.businessNumber<<cBusinessNumberShift | iw.dbId<<cSequenceBits | iw.sequence
	return uID, nil
}

func timeGen() int64 {
	return time.Now().Unix()
}

func timeReGen(last int64) int64 {
	ts := timeGen()
	for {
		if ts <= last {
			ts = timeGen()
		} else {
			break
		}
	}
	return ts
}

func BenchmarkGenNewID_WithLock(b *testing.B) {
	iw := &TestLock{
		businessNumber: 0,
		lastTimeStamp:  0,
		sequence:       0,
		instance:       0,
		algVerId:       0,
		dbCounter:      0,
		lock:           new(sync.Mutex),
		isSub:          false,
		nodesCnt:       1,
	}

	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		iw.NextID()
	}
}

func BenchmarkGenNewID_WithAtomic(b *testing.B) {
	iw := &TestAtomic{
		businessNumber: 0,
		lastTimeStamp:  0,
		sequence:       atomic.NewInt64(0),
		instance:       0,
		algVerId:       0,
		dbCounter:      atomic.NewInt32(0),
		lock:           new(sync.Mutex),
		isSub:          false,
		nodesCnt:       1,
	}

	b.ResetTimer()
	for N := 0; N < b.N; N++ {
		iw.NextID()
	}
}
