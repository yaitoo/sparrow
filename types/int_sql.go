package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (f *Int) Scan(value interface{}) error {
	if value == nil {
		f.int, f.valid = 0, false
		return nil
	}

	var ns sql.NullFloat64

	err := ns.Scan(value)
	if err != nil {
		f.int = int(ns.Float64)
		f.valid = false
		return err
	}

	f.valid = true
	f.int = int(ns.Float64)

	return nil
}

// Value implements the driver.Valuer interface.
func (f Int) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.int, nil
}
