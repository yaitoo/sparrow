package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (f *Float32) Scan(value interface{}) error {
	if value == nil {
		f.float32, f.valid = 0, false
		return nil
	}

	var ns sql.NullFloat64

	err := ns.Scan(value)
	if err != nil {
		f.float32 = float32(ns.Float64)
		f.valid = false
		return err
	}

	f.valid = true
	f.float32 = float32(ns.Float64)

	return nil
}

// Value implements the driver.Valuer interface.
func (f Float32) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.float32, nil
}
