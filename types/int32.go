package types

import (
	"fmt"
)

func NewInt32(val int32) Int32 {
	return Int32{
		int32: val,
		valid: true,
	}
}

//Int32 int32 wrapper
type Int32 struct {
	int32

	valid bool
}

//GetValue return int32 value
func (f *Int32) GetValue() int32 {
	return f.int32
}

//SetValue set int32 value
func (f *Int32) SetValue(i int32) {
	f.int32 = i
	f.valid = true
}

//Ptr return point32 value
func (f *Int32) Ptr() *int32 {
	if f.valid {
		return &f.int32
	}
	return nil
}

//String implements Stringer
func (f Int32) String() string {
	if f.valid {
		return fmt.Sprint(f.int32)
	}

	return ""
}
