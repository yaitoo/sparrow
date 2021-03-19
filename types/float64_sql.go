package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (f *Float64) Scan(value interface{}) error {
	if value == nil {
		f.float64, f.valid = 0, false
		return nil
	}

	var ns sql.NullFloat64

	err := ns.Scan(value)
	if err != nil {
		f.float64 = float64(ns.Float64)
		f.valid = false
		return err
	}

	f.valid = true
	f.float64 = float64(ns.Float64)

	return nil
}

// Value implements the driver.Valuer interface.
func (f Float64) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.float64, nil
}
