package keys

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/yaitoo/sparrow/types"
)

type idGenerator struct {
	ctx        context.Context
	nextIDs    chan int64
	currentID  int64
	weight     int
	high       int64
	low        int64
	timeLayout string
}

// NewHiLo create a IDGenerator instance based on  Hi/Lo algorithm
func NewHiLo(ctx context.Context, low, high int64, weight int, timeFormat string) IDGenerator {
	g := &idGenerator{}
	g.ctx = ctx
	g.currentID = low
	g.high = high
	g.weight = weight
	g.nextIDs = make(chan int64, 1000)
	g.timeLayout = types.FormatLayout(timeFormat)

	go g.startGen(ctx)
	return g
}

func (g *idGenerator) startGen(ctx context.Context) {

	for {
		select {
		case <-g.ctx.Done():
			break
		default:
			g.nextIDs <- g.currentID

			if g.high > 0 && g.currentID >= g.high {
				g.currentID = g.low
			} else {
				g.currentID++
			}

		}
	}
}
func (g *idGenerator) NewID() (string, error) {
	select {
	case n := <-g.nextIDs:
		return strings.ToUpper(types.PadLeft(strconv.FormatInt(n, 36), g.weight, "0")), nil
	case <-time.After(1 * time.Second):
		return "", ErrNoIDs
	}

}

func (g *idGenerator) NewTimeID() (string, error) {
	select {
	case n := <-g.nextIDs:
		t := types.Atoi(time.Now().UTC().Format(g.timeLayout), 0)
		return strings.ToUpper(strconv.FormatInt(int64(t), 36) + "-" + types.PadLeft(strconv.FormatInt(n, 36), g.weight, "0")), nil
	case <-time.After(1 * time.Second):
		return "", ErrNoIDs
	}
}

func (g *idGenerator) NewWith(now time.Time) (string, error) {
	select {
	case n := <-g.nextIDs:
		t := types.Atoi(now.UTC().Format(g.timeLayout), 0)
		return strings.ToUpper(strconv.FormatInt(int64(t), 36) + "-" + types.PadLeft(strconv.FormatInt(n, 36), g.weight, "0")), nil
	case <-time.After(1 * time.Second):
		return "", ErrNoIDs
	}
}
