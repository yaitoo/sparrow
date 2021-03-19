package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (s *String) Scan(value interface{}) error {
	if value == nil {
		s.string, s.valid = "", false
		return nil
	}

	var ns sql.NullString

	err := ns.Scan(value)
	if err != nil {
		s.string = ns.String
		s.valid = false
		return err
	}

	s.valid = true
	s.string = ns.String

	return nil
}

// Value implements the driver.Valuer interface.
func (s String) Value() (driver.Value, error) {
	if !s.valid {
		return nil, nil
	}
	return s.string, nil
}
