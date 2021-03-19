package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner int32erface.
func (f *Int32) Scan(value interface{}) error {
	if value == nil {
		f.int32, f.valid = 0, false
		return nil
	}

	var ns sql.NullFloat64

	err := ns.Scan(value)
	if err != nil {
		f.int32 = int32(ns.Float64)
		f.valid = false
		return err
	}

	f.valid = true
	f.int32 = int32(ns.Float64)

	return nil
}

// Value implements the driver.Valuer int32erface.
func (f Int32) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.int32, nil
}
