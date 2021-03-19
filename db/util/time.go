package util

import (
	"strconv"
	"time"
)

type Clock interface {
	Now() time.Time
	// After(d time.Duration) <-chan time.Time
}

type TimeUtil struct {
	t time.Time
}

func (t TimeUtil) Now() time.Time {
	return time.Now()
}

func NewTimeUtil() TimeUtil {
	return TimeUtil{}
}

func NewTimeUtilWithTime(t time.Time) TimeUtil {
	return TimeUtil{
		t: t,
	}
}

func ConvertWIthTimeZone(t time.Time, tzStr string) (TimeUtil, error) {
	if tzStr == "" {
		return TimeUtil{
			t: t,
		}, nil
	}
	timezone, err := time.LoadLocation(tzStr)
	if err != nil {
		return TimeUtil{t: t}, err
	}

	return TimeUtil{
		t: t.UTC().In(timezone),
	}, nil
}

func (t *TimeUtil) ToDayString() string {
	return strconv.Itoa(t.t.Day())
}

func (t *TimeUtil) ToMonthString() string {
	return strconv.Itoa(int(t.t.Month()))
}

func (t *TimeUtil) ToYearString() string {
	return strconv.Itoa(t.t.Year())
}

func (t *TimeUtil) ToWeekString() string {
	_, week := t.t.ISOWeek()
	return strconv.Itoa(week)
}

func (t *TimeUtil) ToYearDayString() string {
	return strconv.Itoa(t.t.YearDay())
}

func (t *TimeUtil) GetTimeStampByLayout(timeStr, layout string) (int64, error) {
	value, err := time.Parse(layout, timeStr)

	if err != nil {
		return 0, err
	}
	return value.Unix(), nil
}

func (t *TimeUtil) ToUtcTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0).UTC()
}
