package types

import (
	"fmt"
)

func NewInt64(val int64) Int64 {
	return Int64{
		int64: val,
		valid: true,
	}
}

//Int64 int64 wrapper
type Int64 struct {
	int64

	valid bool
}

//GetValue return int64 value
func (f *Int64) GetValue() int64 {
	return f.int64
}

//SetValue set int64 value
func (f *Int64) SetValue(i int64) {
	f.int64 = i
	f.valid = true
}

//Ptr return point64 value
func (f *Int64) Ptr() *int64 {
	if f.valid {
		return &f.int64
	}
	return nil
}

//String implements Stringer
func (f Int64) String() string {
	if f.valid {
		return fmt.Sprint(f.int64)
	}

	return ""
}
