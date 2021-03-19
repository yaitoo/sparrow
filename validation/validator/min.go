package validator

import (
	"fmt"
	"strconv"

	"github.com/yaitoo/sparrow/validation/models"
)

type Min struct {
}

func (m *Min) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	min, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(min, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Minimum is %d", min), false
}

func (m *Min) parseParams(params string) (int, bool) {
	i, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return i, true
}

// IsSatisfied judge whether obj is valid
// not support int64 on 32-bit platform
func (m *Min) isSatisfied(min int, obj interface{}) bool {
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

	return v >= min
}
