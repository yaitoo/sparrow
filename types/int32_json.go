package types

import (
	"encoding/json"
	"strings"
)

// MarshalJSON implements json.Marshaler.
func (f Int32) MarshalJSON() ([]byte, error) {
	if f.valid {
		return json.Marshal(f.int32)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (f *Int32) UnmarshalJSON(data []byte) error {
	var v int32

	if data == nil || len(data) == 0 {
		f.valid = false
		return nil
	}

	s := strings.Trim(string(data), "\"")

	if len(s) == 0 {
		f.valid = false
		return nil
	}

	if s == "null" {
		f.valid = false
		return nil
	}

	err := json.Unmarshal([]byte(s), &v)

	if err == nil {
		f.int32 = v
		f.valid = true
		return nil
	}

	return err

}
