package validator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/yaitoo/sparrow/validation/models"
)

type Range struct {
}

func (m *Range) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	min, max, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(min, max, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Range is %d to %d", min, max), false
}

func (m *Range) parseParams(params string) (int, int, bool) {

	vals := strings.Split(params, ",")
	if len(vals) != 2 {
		return 0, 0, false
	}

	min, err := strconv.Atoi(strings.TrimSpace(vals[0]))
	if err != nil {
		return 0, 0, false
	}

	max, err := strconv.Atoi(strings.TrimSpace(vals[1]))
	if err != nil {
		return 0, 0, false
	}
	return min, max, true
}

// isSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (m *Range) isSatisfied(min, max int, obj interface{}) bool {
	var v int
	switch obj.(type) {
	case int64:
		if wordsize == 32 {
			return false
		}
		v = int(obj.(int64))
	case int:
		v = obj.(int)
	case int32:
		v = int(obj.(int32))
	case int16:
		v = int(obj.(int16))
	case int8:
		v = int(obj.(int8))
	default:
		return false
	}

	return v >= min && v <= max
}
