package keys

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/yaitoo/sparrow/log"
	"github.com/yaitoo/sparrow/types"
)

var logger = log.NewLogger("snowflake")

//ErrNodeOverflow node overflow in snowflake
var ErrNodeOverflow = errors.New("snowflake: node is up to 999")

//ErrCapOverflow  capacity overflow in snowflake
var ErrCapOverflow = errors.New("snowflake: capacity overflow")

//采用类似 Twitter的snowFlake算法生成分布式ID，
//格式采用 时间格式（自定义格式）+节点ID（3位INT）+序号（自定义最高值）
//其中时间可以带入，序号可以跨时间周期复用
type snowFlakeIDGenerator struct {
	sync.Mutex

	current int64
	cap     int64
	weight  int

	nodeID int

	timeZone   int
	timeLayout string

	lastTime string
}

//NewSnowflake 创建雪花ID生成器实例
func NewSnowflake(ctx context.Context, current, cap int64, nodeID int, timeZone int, timeLayout string) (IDGenerator, error) {
	gen := &snowFlakeIDGenerator{}

	if cap <= 0 {
		cap = 10000
	}

	if current >= cap {
		current = 0
	}

	gen.current = current
	gen.cap = cap
	gen.weight = len(strconv.FormatInt(cap, 10))

	if nodeID > 999 {
		return nil, ErrNodeOverflow
	}
	gen.nodeID = nodeID

	gen.timeZone = timeZone
	gen.timeLayout = types.FormatLayout(timeLayout)

	return gen, nil
}

func (gen *snowFlakeIDGenerator) NewID() (string, error) {
	gen.Lock()
	defer gen.Unlock()

	now := types.SwitchTimezone(time.Now(), gen.timeZone)
	currentTime := now.Format(gen.timeLayout)

	current := gen.current
	if current < gen.cap {

		next := currentTime + types.PadLeft(strconv.Itoa(gen.nodeID), 3, "0") + types.PadLeft(strconv.FormatInt(current, 10), gen.weight, "0")

		gen.lastTime = currentTime
		gen.current = current + 1

		return next, nil
	}

	//current >= gen.cap, 序号容量用完

	//时间周期已经更换，则可以直接从0开始
	if currentTime > gen.lastTime {
		gen.lastTime = currentTime
		current = 0
		gen.current = current + 1

		return currentTime + types.PadLeft(strconv.Itoa(gen.nodeID), 3, "0") + types.PadLeft(strconv.FormatInt(current, 10), gen.weight, "0"), nil
	}

	//时间周期未更换，则需要等待下一个时间周期
	for {
		logger.Warnln("snowflake: cap is overflow for ", currentTime)
		now = types.SwitchTimezone(time.Now(), gen.timeZone)
		currentTime = now.Format(gen.timeLayout)
		if currentTime > gen.lastTime {

			gen.lastTime = currentTime

			current = 0
			gen.current = current + 1

			return currentTime + types.PadLeft(strconv.Itoa(gen.nodeID), 3, "0") + types.PadLeft(strconv.FormatInt(current, 10), gen.weight, "0"), nil
		}
		time.Sleep(500 * time.Millisecond)
	}

}

//NewTimeID 指定时间，则可以突破容量
func (gen *snowFlakeIDGenerator) NewTimeID() (string, error) {

	return gen.NewID()
}

//NewWith 指定时间，则可以突破容量
func (gen *snowFlakeIDGenerator) NewWith(t time.Time) (string, error) {
	gen.Lock()
	defer gen.Unlock()

	currentTime := t.Format(gen.timeLayout)
	current := gen.current
	if current < gen.cap {

		gen.current = current + 1

		next := currentTime + types.PadLeft(strconv.Itoa(gen.nodeID), 3, "0") + types.PadLeft(strconv.FormatInt(current, 10), gen.weight, "0")
		return next, nil
	}

	//容量超过，并且时间周期更换， 则直接从0开始计数
	if currentTime > gen.lastTime {
		current = 0

		gen.current = 1
		gen.lastTime = currentTime

		next := currentTime + types.PadLeft(strconv.Itoa(gen.nodeID), 3, "0") + types.PadLeft(strconv.FormatInt(current, 10), gen.weight, "0")
		return next, nil
	}
	logger.Errorln("snowflake: cap overflow, it is up to ", gen.cap)
	return "", ErrCapOverflow
}
