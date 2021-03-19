package types

import (
	"fmt"
)

func NewFloat64(val float64) Float64 {
	return Float64{
		float64: val,
		valid:   true,
	}
}

//Float64 float64 wrapper
type Float64 struct {
	float64

	valid bool
}

//GetValue return Float64 value
func (f *Float64) GetValue() float64 {
	return f.float64
}

//SetValue set Float64 value
func (f *Float64) SetValue(i float64) {
	f.float64 = i
	f.valid = true
}

//Ptr return point value
func (f *Float64) Ptr() *float64 {
	if f.valid {
		return &f.float64
	}
	return nil
}

//String implements Stringer
func (f Float64) String() string {
	if f.valid {
		return fmt.Sprint(f.float64)
	}

	return ""
}
