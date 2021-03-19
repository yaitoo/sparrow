package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/yaitoo/sparrow/validation/models"
)

type MaxSize struct {
}

func (m *MaxSize) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	max, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(max, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Maximum size is %d", max), false
}

func (m *MaxSize) parseParams(params string) (int, bool) {
	i, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return i, true
}

// IsSatisfied judge whether obj is valid
func (m *MaxSize) isSatisfied(max int, obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) <= max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= max
	}
	return false
}
