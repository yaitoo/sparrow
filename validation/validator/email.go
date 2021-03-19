package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type Email struct {
}

func (r *Email) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be a valid email address", false
}

// isSatisfied judge whether obj is valid
func (r *Email) isSatisfied(obj interface{}) bool {
	return emailPattern.MatchString(fmt.Sprintf("%v", obj))
}
