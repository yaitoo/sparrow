package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type AlphaDash struct {
}

func (r *AlphaDash) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be valid alpha or numeric or dash(-_) characters", false
}

// isSatisfied judge whether obj is valid
func (r *AlphaDash) isSatisfied(obj interface{}) bool {
	return !alphaDashPattern.MatchString(fmt.Sprintf("%v", obj))
}
