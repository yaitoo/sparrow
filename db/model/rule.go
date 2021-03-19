package model

import (
	"encoding/json"
)

type Rule struct {
	Tables    []Table    `yaml:tables`
	Databases []Database `yaml:databases`
}

/* func (r *Rule) Validating() bool {
	for _, v := range r.Tables {
		if v.Validating() == false {
			return false
		}
	}
	return true
} */

func (r *Rule) String() string {
	result, _ := json.Marshal(r)
	return string(result)
}

func (r *Rule) IsEmptyObject() bool {
	return len(r.Tables) == 0 //r.buinessTag == ""
}

func (r *Rule) GetBusicessEntity(entityName string) Table {
	for idx := range r.Tables {
		if r.Tables[idx].Name == entityName {
			return r.Tables[idx]
		}
	}
	return Table{}
}

func (r Rule) GetDatabase(dbId int64) Database {
	arrLen := int64(len(r.Databases))
	if arrLen >= dbId {
		return r.Databases[dbId]
	}
	return Database{}
}

func (r *Rule) Validate() bool {
	for i := 0; i < len(r.Databases); i++ {
		if r.Databases[i].Validate() == false {
			return false
		}
	}
	for i := 0; i < len(r.Tables); i++ {
		if r.Tables[i].Validate() == false {
			return false
		}
	}
	return true
}
