package db

//Vars a key/value collection
type Vars struct {
	builder *SQLBuilder
	names   []string
}

//WithNames set names
func (vars *Vars) WithNames(names ...string) *Vars {
	if len(names) == 0 {
		return vars
	}

	if vars.names == nil {
		vars.names = make([]string, 0, len(names))
	}

	for _, name := range names {
		vars.names = append(vars.names, name)
	}

	return vars
}

//WithValues set values
func (vars *Vars) WithValues(values ...interface{}) {
	if len(values) == 0 || len(vars.names) == 0 {
		return
	}

	for i := 0; i < len(values) && i < len(vars.names); i++ {
		vars.builder.Var(vars.names[i], values[i], true)
	}
}
