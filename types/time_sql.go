package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999"
)

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(value interface{}) (err error) {
	if value == nil {
		t.Time, t.valid = time.Time{}, false
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		t.Time, t.valid = v, true
		//return nil
	case []byte:
		t.Time, err = parseTime(string(v), time.UTC)
		t.valid = (err == nil)
		//return nil
	case string:
		t.Time, err = parseTime(v, time.UTC)
		t.valid = (err == nil)
		//return nil
	}

	if t.valid {
		t.Time = SwitchTimezone(t.Time, t.ctx.TimeOffset)
		return nil
	}

	t.valid = false
	return fmt.Errorf("Can't convert %T to time.Time", value)
}

//https://github.com/go-sql-driver/mysql/blob/master/driver.go
func parseTime(str string, loc *time.Location) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("Invalid Time-String: %s", str)
		return
	}

	// Adjust location
	if err == nil && loc != time.UTC {
		y, mo, d := t.Date()
		h, mi, s := t.Clock()
		t, err = time.Date(y, mo, d, h, mi, s, t.Nanosecond(), loc), nil
	}

	return
}

//func parseBinaryTime(num uint64, data []byte, loc *time.Location) (driver.Value, error) {
//	switch num {
//	case 0:
//		return time.Time{}, nil
//	case 4:
//		return time.Date(
//			int(binary.LittleEndian.Uint16(data[:2])), // year
//			time.Month(data[2]),                       // month
//			int(data[3]),                              // day
//			0, 0, 0, 0,
//			loc,
//		), nil
//	case 7:
//		return time.Date(
//			int(binary.LittleEndian.Uint16(data[:2])), // year
//			time.Month(data[2]),                       // month
//			int(data[3]),                              // day
//			int(data[4]),                              // hour
//			int(data[5]),                              // minutes
//			int(data[6]),                              // seconds
//			0,
//			loc,
//		), nil
//	case 11:
//		return time.Date(
//			int(binary.LittleEndian.Uint16(data[:2])), // year
//			time.Month(data[2]),                       // month
//			int(data[3]),                              // day
//			int(data[4]),                              // hour
//			int(data[5]),                              // minutes
//			int(data[6]),                              // seconds
//			int(binary.LittleEndian.Uint32(data[7:11]))*1000, // nanoseconds
//			loc,
//		), nil
//	}
//	return nil, fmt.Errorf("Invalid Time-packet length %d", num)
//}

// Value implements the driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	if !t.valid {
		return nil, nil
	}
	return t.Time.Round(time.Microsecond), nil
}
