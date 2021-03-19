package validator

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

	"github.com/yaitoo/sparrow/validation/models"
)

type Length struct {
}

func (m *Length) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	length, ok := m.parseParams(params)
	if !ok {
		return "invalid params for validator", false
	}

	ok = m.isSatisfied(length, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Required length is %d", length), false
}

func (m *Length) parseParams(params string) (int, bool) {
	i, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return i, true
}

/// IsSatisfied judge whether obj is valid
func (m *Length) isSatisfied(length int, obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return utf8.RuneCountInString(str) == length
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == length
	}
	return false
}
