package types

import (
	"encoding/json"
	"strings"
)

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.valid == false {
		return []byte("null"), nil
	}

	layout := "yyyy-MM-dd HH:mm:ss"
	if IsNotEmpty(t.ctx.TimeJSON) {
		layout = t.ctx.TimeJSON
	}

	tz := SwitchTimezone(t.Time, t.ctx.TimeOffset)

	return json.Marshal(FormatTime(&tz, layout))

}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (t *Time) UnmarshalJSON(data []byte) error {

	if data == nil || len(data) == 0 {
		t.valid = false
		return nil
	}

	s := strings.Trim(string(data), "\"")

	layout := "yyyy-MM-dd HH:mm:ss"
	if IsNotEmpty(t.ctx.TimeJSON) {
		layout = t.ctx.TimeJSON
	}

	tm := ParseTime(s, FormatLayout(layout), nil)

	if tm == nil {
		t.valid = false
	} else {
		t.Time = *tm
		t.valid = true
	}

	return nil
}
