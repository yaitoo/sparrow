package types

import (
	"fmt"
	"strings"
)

// PadLeft x
func PadLeft(s string, lenght int, char string) string {

	for len(s) < lenght {
		s = char + s
	}

	return s
}

//ToString  try convert object to string
func ToString(v interface{}, defaultValue string) string {
	if v == nil {
		return defaultValue
	}

	s, ok := v.(string)
	if ok {
		return s
	}

	return fmt.Sprintf("%s", v)

}

//IsEmpty return true if string has no any valid value, else true
func IsEmpty(s string) bool {
	return len(strings.Trim(s, " ")) == 0
}

//IsNotEmpty return true if string has valid value, else true
func IsNotEmpty(s string) bool {
	return IsEmpty(s) == false
}

//LastString return last num char as string
func LastString(s string, num int) string {
	if num > 0 {

		if len(s) > num {
			return s[len(s)-num:]
		}

		return s
	}
	return ""
}
