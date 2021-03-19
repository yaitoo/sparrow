// +---+-----------+------------------+-------------------------------+----------+-----------+---------------+-------------------+
// | 0 | algVer(3) | InstanceId(3) | timestamp(second)30  | tags (3) | db id(4)  | sequence(20) |
// +---+-----------+------------------+-------------------------------+---------------------------+-----------+-------------------+

package shardingId

import (
	"errors"
	"math"
	"strconv"
	"sync"
	"time"

	"go.uber.org/atomic"

	configmanager "github.com/yaitoo/sparrow/db/model"
)

// todo : every 3hr to reset  db counter for genrerator.

const nodeIdRegex = `.\.*(\d{2})$`

var TimeNow time.Time

type errorEnum int

const (
	invalidWorkerID     errorEnum = 0
	clockMovedBackwards errorEnum = 1
)

func (ee errorEnum) String() string {
	return [...]string{"ID should between 0 and %d", "ClockMovedBackwards"}[ee]
}

const (

	// cEpoch : Monday, May 20, 2019 12:00:00 AM， 时间戳起始时间点
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

// ID generator structure
type IdGenerator struct {
	algVerId       int64
	instance       int64
	businessNumber int64
	dbId           int64
	lastTimeStamp  int64
	sequence       *atomic.Int64

	table     string
	tag       string
	entity    configmanager.Table
	lock      *sync.Mutex
	dbCounter *atomic.Int32 // [ 0 ,  len(db nodes)]
	nodesCnt  int
	isSub     bool
	//	vc         *configmanager.VersionRules
}

// ID structure
type IdStruct struct {
	//Id Id值
	Id int64
	//BusinessNumber 根据tag计算出来业务hash值
	BusinessNumber int64
	//DbId 分库ID
	DbId int64
	//Sequence 当前顺序号
	Sequence int64
	//InstanceId 分库实例
	InstanceId int64
	//AlgVer 分库分表Version
	AlgVer int64
	//TimeStamp 时间戳
	TimeStamp int64
	Time      time.Time
}

/* func (i *IdStruct) IsNullObject() bool {
	return i.Id == 0
} */

func (i IdStruct) IsNullObject() bool {
	return i.TimeStamp == 0 || i.TimeStamp < 10000000 || i.DbId > 15 || i.BusinessNumber > 7 || i.AlgVer > 7 || i.InstanceId > 7
}

// NewIDWorker : init id generator for specific business entity
func NewIdGenerator(table string, tag string, vc configmanager.Version, instanceId int64) (iw *IdGenerator, err error) {
	_tableName := getTableName(tag)

	iw = &IdGenerator{
		table: table,
		tag:   _tableName,
	}

	//  get config of this entity
	rule := vc.GetBusicessRule(table)
	if rule.Tables == nil {
		iw.nodesCnt = 1
	} else {
		iw.nodesCnt = len(rule.Databases)
	}

	// init ther remaind properties
	iw.entity = rule.GetBusicessEntity(table)
	iw.businessNumber = int64(iw.entity.GetTable(_tableName).ID)
	iw.lastTimeStamp = 0
	iw.sequence = atomic.NewInt64(0)
	iw.instance = instanceId
	iw.algVerId = vc.GetVer()
	iw.dbCounter = atomic.NewInt32(0)
	iw.lock = new(sync.Mutex)
	iw.isSub = false
	return iw, nil
}

func NewSubId(table string, tag string, vc configmanager.Version, parentId int64) (iw *IdGenerator, err error) {
	_idObject := ParseId(parentId)

	_tableName := getTableName(tag)
	iw = &IdGenerator{
		table: table,
		tag:   _tableName,
	}

	//  get config of this entity
	rule := vc.GetBusicessRule(table)
	iw.nodesCnt = len(rule.Databases)
	// init ther remaind properties
	iw.entity = rule.GetBusicessEntity(table)
	iw.businessNumber = int64(iw.entity.GetTable(_tableName).ID)
	iw.lastTimeStamp = 0
	iw.sequence = atomic.NewInt64(0)
	iw.instance = _idObject.InstanceId
	iw.algVerId = vc.GetVer()
	iw.dbCounter = atomic.NewInt32(0)
	iw.lock = new(sync.Mutex)
	iw.dbId = _idObject.DbId
	iw.isSub = true
	return iw, nil
}

func getTableName(tableNm string) string {
	if tableNm == "" {
		return "*"
	}
	return tableNm
}

// NextID : get next ID
func (iw *IdGenerator) NextID() (uID int64, err error) {
	//iw.lock.Lock()
	//defer iw.lock.Unlock()
	/* 	if err != nil {
		return 0, err
	} */
	if iw.isSub == false {
		iw.lock.Lock()
		aDBCounter := iw.dbCounter.Load()
		iw.dbId = int64((aDBCounter) % int32(iw.nodesCnt))
		if aDBCounter == math.MaxInt32 {
			iw.dbCounter.Store(0)
		}
		iw.dbCounter.Inc()
		iw.lock.Unlock()
	}

	// 產生新的時間
	ts := timeGen()
	iw.lock.Lock()
	nextSequence := iw.sequence.Inc()
	// 超過 20 位的二進位，就歸零等待下一秒到來。
	if nextSequence >= 1048575 {
		// block 直到下一秒到來
		ts = timeReGen(ts)
		// 歸零
		iw.sequence.Store(0)
		nextSequence = 0
	}
	iw.lock.Unlock()
	//if ts == iw.lastTimeStamp {
	//	iw.sequence = (iw.sequence + 1) & cSequenceMask
	//	if iw.sequence == 0 {
	//		ts = timeReGen(ts)
	//	}
	//} else {
	//	iw.sequence = 0
	//}

	// 如果系統的時間被往前調整，就會報錯。
	if ts < iw.lastTimeStamp {
		err = errors.New(clockMovedBackwards.String())
		return 0, err
	}

	iw.lastTimeStamp = ts
	uID = iw.algVerId<<cAlgVersionShift | iw.instance<<cInstanceShift | (ts-cEpoch)<<cTimeStampShift | iw.businessNumber<<cBusinessNumberShift | iw.dbId<<cSequenceBits | nextSequence
	return uID, nil
}

func timeGen() int64 {
	if TimeNow.IsZero() {
		return time.Now().Unix()
	}
	return TimeNow.Unix()
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

// ParseId : parse id to IdStruct
func ParseId(id int64) IdStruct {
	result := IdStruct{
		Id:             id,
		Sequence:       id & cSequenceMask,
		DbId:           (id >> cSequenceBits) & cMaxDb,
		BusinessNumber: (id >> cBusinessNumberShift) & cMaxDBusinessNumber,
		TimeStamp:      (id >> cTimeStampShift) & cMaxTimestamp,
		InstanceId:     (id >> cInstanceShift) & cMaxInstance,
		AlgVer:         (id >> cAlgVersionShift) & cMaxAlgVer,
	}
	result.Time = time.Unix(result.TimeStamp+cEpoch, 0).UTC()
	return result
}

// ParseId : parse id to IdStruct
func ParseStringId(id string) IdStruct {
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return IdStruct{}
	}

	result := IdStruct{
		Id:             idInt64,
		Sequence:       idInt64 & cSequenceMask,
		DbId:           (idInt64 >> cSequenceBits) & cMaxDb,
		BusinessNumber: (idInt64 >> cBusinessNumberShift) & cMaxDBusinessNumber,
		TimeStamp:      (idInt64 >> cTimeStampShift) & cMaxTimestamp,
		InstanceId:     (idInt64 >> cInstanceShift) & cMaxInstance,
		AlgVer:         (idInt64 >> cAlgVersionShift) & cMaxAlgVer,
	}
	result.Time = time.Unix(result.TimeStamp+cEpoch, 0).UTC()
	return result
}
