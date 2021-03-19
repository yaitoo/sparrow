package validator

import (
	"fmt"
	"regexp"

	"github.com/yaitoo/sparrow/validation/models"
)

type Match struct {
}

func (r *Match) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(params, value)
	if ok {
		return "", true
	}

	return fmt.Sprintf("Must match %s", params), false
}

// isSatisfied judge whether obj is valid
func (r *Match) isSatisfied(pattern string, obj interface{}) bool {
	matched, _ := regexp.MatchString(pattern, fmt.Sprintf("%v", obj))
	return matched
}
