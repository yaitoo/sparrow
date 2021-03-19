package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type Base64 struct {
}

func (r *Base64) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid base64 characters", false
}

// isSatisfied judge whether obj is valid
func (r *Base64) isSatisfied(obj interface{}) bool {
	return base64Pattern.MatchString(fmt.Sprintf("%v", obj))
}
