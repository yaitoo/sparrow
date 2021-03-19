package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type Mobile struct {
}

func (r *Mobile) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid mobile number", false
}

// isSatisfied judge whether obj is valid
func (r *Mobile) isSatisfied(obj interface{}) bool {
	return mobilePattern.MatchString(fmt.Sprintf("%v", obj))
}
