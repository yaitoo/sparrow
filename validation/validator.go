package validation

import (
	"strings"
	"sync"

	"github.com/yaitoo/sparrow/validation/models"
)

var (
	mutex      sync.RWMutex
	validators = make(map[string]Validator)
)

//Validator  implement this interface in order to register your custom validator
type Validator interface {
	Validate(form models.Form, params, key string, value interface{}) (string, bool)
}

// Register a validator with name
func Register(name string, validator Validator) {

	if validator == nil {
		return
	}

	name = strings.ToLower(strings.TrimSpace(name))
	if len(name) == 0 {
		panic("Validator name is requried")
	}

	mutex.Lock()
	defer mutex.Unlock()

	if _, dup := validators[name]; dup {
		logger.Warnln(name, " is duplicated")
	}

	validators[name] = validator
}

//getValidator the validator with the register name
func getValidator(rule string) (Validator, string) {

	if len(rule) == 0 {

		return nil, ""
	}

	idx := strings.Index(rule, ":")

	if idx == 0 {
		return nil, ""
	}

	name, params := rule, ""
	if idx > 0 {
		name = rule[:idx]
		params = rule[idx+1:]
	}

	name = strings.ToLower(strings.TrimSpace(name))

	mutex.RLock()
	defer mutex.RUnlock()

	validator, _ := validators[name]
	return validator, params
}
