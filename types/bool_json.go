package types

import (
	"encoding/json"
	"strings"
)

// MarshalJSON implements json.Marshaler.
func (b Boolean) MarshalJSON() ([]byte, error) {
	if b.valid {
		return json.Marshal(b.bool)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (b *Boolean) UnmarshalJSON(data []byte) error {

	if data == nil || len(data) == 0 {
		b.valid = false
		return nil
	}

	s := strings.Trim(string(data), "\"")

	if len(s) == 0 {
		b.valid = false
		return nil
	}

	if s == "null" {
		b.valid = false
		return nil
	}

	// if strings.ToLower(s) == "true" {
	// 	b.bool = true
	// }

	// b.valid = true

	var v bool
	err := json.Unmarshal([]byte(s), &v)

	if err == nil {
		b.bool = v
		b.valid = true
		return nil
	}
	return err

}
