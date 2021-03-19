package validation

import (
	"context"
	"strings"
	"sync"

	"github.com/yaitoo/sparrow/validation/models"
)

//Context validation context
type Context struct {
	sync.RWMutex
	ctx         context.Context
	formsData   map[string]map[string]interface{} // map[section](map[key]value)
	errors      []string
	isValidated bool
	forms       map[string]models.Form
}

// Group help to add keys and values to a group quickly
type Group struct {
	c      *Context
	group  string
	keys   []string
	values []interface{}
}

//NewContext create a new validation context
func NewContext(ctx context.Context, options ...Option) *Context {
	c := &Context{ctx: ctx, formsData: make(map[string]map[string]interface{})}

	for _, option := range options {
		option(c)
	}

	if c.forms == nil {
		c.forms = forms
	}

	return c
}

//Validate validate a parameter
func (c *Context) Validate(formName, key string, value interface{}) *Context {
	c.updateValue(formName, key, value)
	return c
}

//ValidateGroup return a validation group
func (c *Context) ValidateGroup(formName string) *Group {
	return &Group{c: c, group: formName}
}

//Error return all validation errors
func (c *Context) Error() string {
	c.RLock()
	defer c.RUnlock()

	return strings.Join(c.errors, ",")
}

//Errors return all validation errors
func (c *Context) Errors() []string {
	c.RLock()
	defer c.RUnlock()

	return c.errors
}

// Clear reset context's status
func (c *Context) Clear() {
	c.Lock()
	defer c.Unlock()

	c.formsData = make(map[string]map[string]interface{})
	c.errors = make([]string, 0)
	c.isValidated = false
}

// getLang get lang from context
func (c *Context) getLang() string {
	lang, ok := c.ctx.Value(defaultLangCtxKey).(string)
	if !ok {
		return defaultLang
	}
	if lang == "" {
		return defaultLang
	}

	return lang
}

//IsValid return validataion result
func (c *Context) IsValid() bool {
	c.Lock()
	defer c.Unlock()

	lang := c.getLang()
	c.errors = make([]string, 0)

	for formName, fields := range c.formsData {
		form, ok := c.forms[formName]
		if !ok {
			logger.Printf("form[%s] is missing\n", formName)
			continue
		}

		for key, value := range fields {
			rules, ok := form[key]
			if !ok || len(rules) == 0 {
				logger.Printf("field[%s.%s] is missing\n", formName, key)
				continue
			}
			for _, rule := range rules {
				validator, validatorParams := getValidator(rule.Rule)
				if validator == nil {
					logger.Warnln("validator[%s] isn't registered", rule.Rule)
					continue
				}

				if err, passed := validator.Validate(form, validatorParams, formName+"."+key, value); passed == false {

					if msg, ok := rule.Message[lang]; ok {
						err = msg
					}

					c.errors = append(c.errors, err)
				}

			}
		}
	}

	return len(c.errors) == 0
}

func (c *Context) updateForm(formName string, keys []string, values []interface{}) {
	c.Lock()
	defer c.Unlock()

	form, ok := c.formsData[formName]
	if !ok {
		form = make(map[string]interface{})
		c.formsData[formName] = form
	}

	if len(keys) != len(values) {
		logger.Errorf("keys length %d not equals to values length %d", len(keys), len(values))
		return
	}

	for idx, key := range keys {
		form[key] = values[idx]
	}
}

func (c *Context) updateValue(formName string, key string, value interface{}) {
	c.updateForm(formName, []string{key}, []interface{}{value})
}

//WithNames push parameters
func (g *Group) WithNames(keys ...string) *Group {
	g.keys = keys
	return g
}

//WithValues push values
func (g *Group) WithValues(values ...interface{}) *Context {
	if len(g.keys) != len(values) {
		logger.Errorf("keys length %d not equals to values length %d", len(g.keys), len(values))
		return g.c
	}
	g.values = values
	g.c.updateForm(g.group, g.keys, g.values)
	return g.c
}
