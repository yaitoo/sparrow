package types

func NewString(val string) String {
	return String{
		string: val,
		valid:  true,
	}
}

//String string wrapper
type String struct {
	string

	valid bool
}

//GetValue return string value
func (s *String) GetValue() string {
	return s.string
}

//SetValue set string value
func (s *String) SetValue(v string) {
	s.string = v
	s.valid = true
}

//HasValue return true if it is no empty
func (s *String) HasValue() bool {
	return IsNotEmpty(s.string)
}

//Ptr return point value
func (s *String) Ptr() *string {
	if s.valid {
		return &s.string
	}
	return nil
}

//String implements Stringer
func (s String) String() string {
	return s.string
}

// //SetContext implements Object.SetContext
// func (s *String) SetContext(ctx context.Context) {
// 	s.ctx = ctx
// }

//IsNull implements Object.IsNull
func (s *String) IsNull() bool {
	return s.valid == false
}
