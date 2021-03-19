package validator

import (
	"fmt"
	"strconv"

	"github.com/yaitoo/sparrow/validation/models"
)

type Max struct {
}

func (m *Max) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	max, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(max, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Maximum is %d", max), false
}

func (m *Max) parseParams(params string) (int, bool) {
	i, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return i, true
}

// isSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (m *Max) isSatisfied(max int, obj interface{}) bool {
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

	return v <= max
}
