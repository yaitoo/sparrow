package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/yaitoo/sparrow/validation/models"
)

type MinSize struct {
}

func (m *MinSize) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	min, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(min, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Minimum size is %d", min), false
}

func (m *MinSize) parseParams(params string) (int, bool) {
	i, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return i, true
}

/// IsSatisfied judge whether obj is valid
func (m *MinSize) isSatisfied(min int, obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) >= min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= min
	}
	return false
}
