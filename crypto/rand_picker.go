package crypto

import (
	"crypto/rand"
	"errors"
	"math/big"
)

//ErrNoElementItem  列表裡沒有任何可用對象
var ErrNoElementItem = errors.New("crypto/rand: no element item to pick")

//RandPicker 數組隨機項抽取器
type RandPicker interface {
	//Add 增加一個選項
	Add(items ...interface{})

	//HasNext 是否還有下一個對象
	HasNext() bool

	//隨機挑選下一個值，並且從待選數組裡面移除
	Next() (interface{}, error)

	//隨機挑選一個值，但是不做移除
	Rand() (interface{}, error)
}

//NewRandPicker 創建數組隨機抽取器實例
func NewRandPicker(max int) RandPicker {
	return &randPicker{
		items: make([]interface{}, 0, max),
	}
}

type randPicker struct {
	items []interface{}
}

func (rp *randPicker) Add(items ...interface{}) {
	rp.items = append(rp.items, items...)
}

func (rp *randPicker) HasNext() bool {
	return len(rp.items) > 0
}

func (rp *randPicker) Next() (interface{}, error) {
	n := len(rp.items)
	if n > 0 {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			return -1, err
		}

		nextIndex := int(i.Int64())
		next := rp.items[nextIndex]

		items := append(rp.items[:nextIndex], rp.items[nextIndex+1:]...)

		rp.items = items

		return next, nil
	}

	return 0, ErrNoElementItem
}

func (rp *randPicker) Rand() (interface{}, error) {
	n := len(rp.items)
	if n > 0 {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			return -1, err
		}

		nextIndex := int(i.Int64())
		next := rp.items[nextIndex]

		return next, nil
	}

	return 0, ErrNoElementItem
}

//IntsRandPicker 隨機選取一個int值，直到全部int取完為止
type IntsRandPicker struct {
	items []int
	//	Next() int
}

//NewIntsRandPicker 生成int隨機
func NewIntsRandPicker(min, max int) *IntsRandPicker {

	var from, to int

	if max > min {
		from = min
		to = max
	} else {
		from = max
		to = min
	}

	items := make([]int, 0, from-to)

	for i := from; i < to; i++ {
		items = append(items, i)
	}

	return &IntsRandPicker{
		items: items,
	}
}

//NewIntsRandPickerItems 生成int隨機
func NewIntsRandPickerItems(items ...int) *IntsRandPicker {
	return &IntsRandPicker{
		items: items,
	}
}

//HasNext 是否還有下一個元素
func (irp *IntsRandPicker) HasNext() bool {
	return len(irp.items) > 0
}

//Next 隨機挑選下一個值，並且從待選數組裡面移除，非線程安全
func (irp *IntsRandPicker) Next() (int, error) {
	n := len(irp.items)
	if n > 0 {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			return -1, err
		}

		nextIndex := int(i.Int64())
		next := irp.items[nextIndex]

		items := append(irp.items[:nextIndex], irp.items[nextIndex+1:]...)

		irp.items = items

		return next, nil
	}

	return 0, ErrNoElementItem
}

//Rand 隨機挑選下一個值，但是不從數組裡移除，非線程安全
func (irp *IntsRandPicker) Rand() (int, error) {
	n := len(irp.items)
	if n > 0 {
		i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
		if err != nil {
			return -1, err
		}

		nextIndex := int(i.Int64())
		next := irp.items[nextIndex]

		return next, nil
	}

	return 0, ErrNoElementItem
}
