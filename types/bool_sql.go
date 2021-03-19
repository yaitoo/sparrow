package types

import (
	"database/sql/driver"
)

// Scan implements the sql.Scanner interface.
func (b *Boolean) Scan(value interface{}) error {
	if value == nil {
		b.bool, b.valid = false, false
		return nil
	}

	v, ok := value.([]byte)

	if ok {
		b.valid = len(v) > 0
		if b.valid {
			b.bool = v[0] == 1
		}
		return nil
	}

	b.valid = false
	b.bool = false

	return nil
}

// Value implements the driver.Valuer interface.
func (b Boolean) Value() (driver.Value, error) {
	if !b.valid {
		return nil, nil
	}

	if b.bool {
		return []byte{1}, nil
	}

	return []byte{0}, nil

	//return b.bool, nil
}
