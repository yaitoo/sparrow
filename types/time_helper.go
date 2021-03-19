package types

import (
	"strings"
	"time"
)

const (

	//TimeLongMonth long format of month
	TimeLongMonth = "January"
	//TimeMonth format of month
	TimeMonth = "Jan"
	//TimeNumMonth number format of month
	TimeNumMonth = "1"
	//TimeZeroMonth zero format of  month
	TimeZeroMonth = "01"
	//TimeLongWeekDay long format of weekday
	TimeLongWeekDay = "Monday"
	//TimeWeekDay format of weekday
	TimeWeekDay = "Mon"
	//TimeDay format of day
	TimeDay = "2"

	//TimeZeroDay  zero format of day
	TimeZeroDay = "02"
	//TimeHour24 24 hours format of hour
	TimeHour24 = "15"
	//TimeHour12 12 hours format of hour
	TimeHour12 = "3"
	//TimeZeroHour12  12 hours zero format of hour
	TimeZeroHour12 = "03"
	//TimeMinute format of minute
	TimeMinute = "4"
	//TimeZeroMinute zero format of minute
	TimeZeroMinute = "04"
	//TimeSecond format of second
	TimeSecond = "5"
	//TimeZeroSecond zero format of second
	TimeZeroSecond = "05"
	//TimeLongYear long format of year
	TimeLongYear = "2006"
	//TimeYear format of year
	TimeYear = "06"
	//TimePM format of PM
	TimePM = "PM"
	//Timepm format of pm
	Timepm = "pm"
	//TimeTZ  MST
	TimeTZ = "MST"
	//TimeISO8601TZ ISO8601TZ
	TimeISO8601TZ = "Z0700" // prints Z for UTC
	//TimeISO8601SecondsTZ ISO8601SecondsTZ
	TimeISO8601SecondsTZ = "Z070000"
	//TimeISO8601ShortTZ ISO8601ShortTZ
	TimeISO8601ShortTZ = "Z07"
	//TimeISO8601ColonTZ ISO8601ColonTZ
	TimeISO8601ColonTZ = "Z07:00" // prints Z for UTC
	//TimeISO8601ColonSecondsTZ ISO8601ColonSecondsTZ
	TimeISO8601ColonSecondsTZ = "Z07:00:00"
	//TimeNumTZ NumTZ
	TimeNumTZ = "-0700" // always numeric
	//TimeNumSecondsTz NumSecondsTz
	TimeNumSecondsTz = "-070000"
	//TimeNumShortTZ NumShortTZ
	TimeNumShortTZ = "-07" // always numeric
	//TimeNumColonTZ NumColonTZ
	TimeNumColonTZ = "-07:00" // always numeric
	//TimeNumColonSecondsTZ NumColonSecondsTZ
	TimeNumColonSecondsTZ = "-07:00:00"
)

const (
	//TimeFormatDateTime yyyyMMDDHHmmSS
	TimeFormatDateTime = TimeLongYear + TimeZeroMonth + TimeZeroDay + TimeHour24 + TimeZeroMinute + TimeZeroSecond
	//TimeFormatDateTimeWithDash yyyy-MM-DD HH:mm:ss
	TimeFormatDateTimeWithDash = TimeLongYear + "-" + TimeZeroMonth + "-" + TimeZeroDay + " " + TimeHour24 + ":" + TimeZeroMinute + ":" + TimeZeroSecond
)

//FormatLayout convert RFC layout to golang magic time number
func FormatLayout(layout string) string {
	f := strings.Replace(layout, "yyyy", TimeLongYear, -1)
	f = strings.Replace(f, "yy", TimeYear, -1)
	f = strings.Replace(f, "MM", TimeZeroMonth, -1)
	f = strings.Replace(f, "M", TimeNumMonth, -1)
	f = strings.Replace(f, "dd", TimeZeroDay, -1)
	f = strings.Replace(f, "d", TimeDay, -1)

	f = strings.Replace(f, "HH", TimeHour24, -1)
	f = strings.Replace(f, "hh", TimeZeroHour12, -1)
	f = strings.Replace(f, "h", TimeHour12, -1)
	f = strings.Replace(f, "mm", TimeZeroMinute, -1)
	f = strings.Replace(f, "m", TimeMinute, -1)
	f = strings.Replace(f, "ss", TimeZeroSecond, -1)
	f = strings.Replace(f, "s", TimeSecond, -1)

	return f
}

//FormatTime format time with special format, eg. YYYY-MM-dd HH:mm:ss
func FormatTime(t *time.Time, format string) string {
	if t == nil {
		return ""
	}

	return t.Format(FormatLayout(format))
}

//SwitchTimezone convert timezone
func SwitchTimezone(t time.Time, offset int) time.Time {
	return t.In(time.FixedZone("", offset*60*60))
}

//ParseTime try to parse string with specified layout, if not, return d
func ParseTime(v string, layout string, d *time.Time) *time.Time {
	if t, err := time.Parse(layout, v); err == nil {
		return &t
	}
	return d
}

//Now 获取指定时区的当前时间
func Now(offset int) time.Time {
	return time.Now().In(time.FixedZone("", offset*60*60))
}

//Today 获取指定时区的今天零点
func Today(offset int) time.Time {
	now := Now(offset)
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
