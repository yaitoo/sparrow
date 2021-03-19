package types

import (
	"fmt"
)

func NewFloat32(val float32) Float32 {
	return Float32{
		float32: val,
		valid:   true,
	}
}

//Float32 float32 wrapper
type Float32 struct {
	float32

	valid bool
}

//GetValue return float32 value
func (f *Float32) GetValue() float32 {
	return f.float32
}

//SetValue set float32 value
func (f *Float32) SetValue(i float32) {
	f.float32 = i
	f.valid = true
}

//Ptr return point value
func (f *Float32) Ptr() *float32 {
	if f.valid {
		return &f.float32
	}
	return nil
}

//String implements Stringer
func (f Float32) String() string {
	if f.valid {
		return fmt.Sprint(f.float32)
	}

	return ""
}
