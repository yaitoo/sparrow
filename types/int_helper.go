package types

import (
	"fmt"
	"strconv"
	"strings"
)

//Atoi convert string to int
func Atoi(v string, d int) int {

	s := strings.Trim(v, "\"")

	if len(s) == 0 {
		return d
	}
	s = strings.TrimLeft(s, "0")
	if len(s) == 0 { //v == "0"
		return 0
	}

	if i, err := strconv.Atoi(s); err == nil {
		return i
	}

	return d
}

//Atoi64 convert string to int64
func Atoi64(v string, d int64) int64 {

	s := strings.Trim(v, "\"")

	if len(s) == 0 {
		return d
	}
	s = strings.TrimLeft(s, "0")
	if len(s) == 0 { //v == "0"
		return 0
	}

	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	return d
}

//ToInt try convert object to int
func ToInt(v interface{}, defaultValue int) int {
	if v == nil {
		return defaultValue
	}

	s, ok := v.(int)
	if ok {
		return s
	}

	return Atoi(fmt.Sprintf("%s", v), defaultValue)

}
