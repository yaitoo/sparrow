package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type Phone struct {
}

func (r *Phone) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid telephone or mobile phone number", false
}

// isSatisfied judge whether obj is valid
func (r *Phone) isSatisfied(obj interface{}) bool {
	s := fmt.Sprintf("%v", obj)
	return telPattern.MatchString(s) || mobilePattern.MatchString(s)
}
