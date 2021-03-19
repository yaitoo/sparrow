package types

import (
	"context"
	"time"
	//"github.com/yaitoo/sparrow/micro"
)

func NewTime(val time.Time) Time {
	return Time{
		Time:  val,
		valid: true,
	}
}

//Time time wrapper
type Time struct {
	time.Time

	valid bool

	ctx Context
}

//GetValue return string value
func (t *Time) GetValue() time.Time {
	return t.Time
}

//SetValue set string value
func (t *Time) SetValue(v time.Time) {
	t.Time = v
	t.valid = true
}

//Ptr return point value
func (t *Time) Ptr() *time.Time {
	if t.valid {
		return &t.Time
	}
	return nil
}

//String implements Stringer
func (t Time) String() string {
	if t.valid {

		tm := SwitchTimezone(t.Time, t.ctx.TimeOffset)

		if IsEmpty(t.ctx.TimeLayout) {
			return tm.String()
		}

		return FormatTime(&tm, t.ctx.TimeLayout)

	}

	return ""
}

//SetContext implements Object.SetContext
func (t *Time) SetContext(ctx context.Context) {
	t.ctx = FromContext(ctx)
}

//IsNull implements Object.IsNull
func (t *Time) IsNull() bool {
	return t.valid == false
}
