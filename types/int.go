package types

import (
	"fmt"
)

func NewInt(val int) Int {
	return Int{
		int:   val,
		valid: true,
	}
}

//Int int wrapper
type Int struct {
	int

	valid bool
}

//GetValue return int value
func (f *Int) GetValue() int {
	return f.int
}

//SetValue set int value
func (f *Int) SetValue(i int) {
	f.int = i
	f.valid = true
}

//Ptr return point value
func (f *Int) Ptr() *int {
	if f.valid {
		return &f.int
	}
	return nil
}

//String implements Stringer
func (f Int) String() string {
	if f.valid {
		return fmt.Sprint(f.int)
	}

	return ""
}
