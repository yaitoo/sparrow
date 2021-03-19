package types

import "fmt"

func NewBoolean(val bool) Boolean {
	return Boolean{
		bool:  val,
		valid: true,
	}
}

//Boolean bool wrapper
type Boolean struct {
	bool

	valid bool
}

//GetValue return bool value
func (f *Boolean) GetValue() bool {
	return f.bool
}

//SetValue set bool value
func (f *Boolean) SetValue(i bool) {
	f.bool = i
	f.valid = true
}

//Ptr return pobool value
func (f *Boolean) Ptr() *bool {
	if f.valid {
		return &f.bool
	}
	return nil
}

//String implements Stringer
func (f Boolean) String() string {
	if f.valid {
		return fmt.Sprint(f.bool)
	}

	return ""
}
