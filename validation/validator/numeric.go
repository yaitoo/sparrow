package validator

import "github.com/yaitoo/sparrow/validation/models"

type Numeric struct {
}

func (r *Numeric) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid numeric characters", false
}

// isSatisfied judge whether obj is valid
func (r *Numeric) isSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		for _, v := range str {
			if '9' < v || v < '0' {
				return false
			}
		}
		return true
	}
	return false
}
