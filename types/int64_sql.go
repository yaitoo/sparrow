package types

import (
	"database/sql"
	"database/sql/driver"
)

// Scan implements the sql.Scanner int64erface.
func (f *Int64) Scan(value interface{}) error {
	if value == nil {
		f.int64, f.valid = 0, false
		return nil
	}

	var ns sql.NullFloat64

	err := ns.Scan(value)
	if err != nil {
		f.int64 = int64(ns.Float64)
		f.valid = false
		return err
	}

	f.valid = true
	f.int64 = int64(ns.Float64)

	return nil
}

// Value implements the driver.Valuer int64erface.
func (f Int64) Value() (driver.Value, error) {
	if !f.valid {
		return nil, nil
	}
	return f.int64, nil
}
