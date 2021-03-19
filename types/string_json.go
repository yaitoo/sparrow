package types

import (
	"encoding/json"
	"strings"
)

// MarshalJSON implements json.Marshaler.
func (s String) MarshalJSON() ([]byte, error) {
	if s.valid == false {
		return []byte("null"), nil
	}

	return json.Marshal(s.string)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (s *String) UnmarshalJSON(data []byte) error {

	if data == nil || len(data) == 0 {
		s.valid = false
		return nil
	}

	src := string(data)

	if src == "null" {
		s.valid = false
		return nil
	}

	if strings.HasPrefix(src, "\"") == false && strings.HasSuffix(src, "\"") == false {
		s.string = src
		s.valid = true
		return nil
	}

	var v string
	err := json.Unmarshal(data, &v)

	if err == nil {
		s.string = v
		s.valid = true
		return nil
	}

	return err

}
