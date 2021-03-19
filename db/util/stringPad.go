package util

import (
	"strconv"
	"strings"
)

func times(str string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(str, n)
}

// Left left-pads the string with pad up to len runes
// len may be exceeded if
func Left(str string, length int, pad string) string {
	return times(pad, length-len(str)) + str
}

// Right right-pads the string with pad up to len runes
func Right(str string, length int, pad string) string {
	return str + times(pad, length-len(str))
}

// Left left-pads the int64 with pad up to len runes
func Int64Left(value int64, length int, pad string) string {
	return times(pad, length-len(string(value))) + strconv.FormatInt(value, 10)
}

func IntLeft(value int, length int, pad string) string {
	return times(pad, length-len(string(value))) + strconv.Itoa(value)
}
