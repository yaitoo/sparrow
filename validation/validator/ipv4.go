package validator

import (
	"fmt"

	"github.com/yaitoo/sparrow/validation/models"
)

type IPv4 struct {
}

func (r *IPv4) Validate(form models.Form, params string, key string, value interface{}) (string, bool) {
	ok := r.isSatisfied(value)
	if ok {
		return "", true
	}

	return "Must be a valid ip address", false
}

// isSatisfied judge whether obj is valid
func (r *IPv4) isSatisfied(obj interface{}) bool {
	return ipPattern.MatchString(fmt.Sprintf("%v", obj))
}
