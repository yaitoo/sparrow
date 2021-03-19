package validator

import "github.com/yaitoo/sparrow/validation/models"

type Alpha struct {
}

func (r *Alpha) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid alpha characters", false
}

// isSatisfied judge whether obj is valid
func (r *Alpha) isSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if ('Z' < v || v < 'A') && ('z' < v || v < 'a') {
				return false
			}
		}
		return true
	}
	return false
}
