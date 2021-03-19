package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type TEL struct {
}

func (r *TEL) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid telephone number", false
}

// isSatisfied judge whether obj is valid
func (r *TEL) isSatisfied(obj interface{}) bool {
	return telPattern.MatchString(fmt.Sprintf("%v", obj))
}
