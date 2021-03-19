package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type ZipCode struct {
}

func (r *ZipCode) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid zipcode", false
}

// isSatisfied judge whether obj is valid
func (r *ZipCode) isSatisfied(obj interface{}) bool {
	return zipCodePattern.MatchString(fmt.Sprintf("%v", obj))
}
