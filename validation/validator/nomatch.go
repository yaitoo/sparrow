package validator

import (
	"fmt"
	"regexp"

	"github.com/yaitoo/sparrow/validation/models"
)

type NoMatch struct {
}

func (r *NoMatch) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(params, value)
	if ok {
		return "", true
	}

	return "Must be valid numeric characters", false
}

// isSatisfied judge whether obj is valid
func (r *NoMatch) isSatisfied(pattern string, obj interface{}) bool {
	matched, _ := regexp.MatchString(pattern, fmt.Sprintf("%v", obj))
	return !matched
}
