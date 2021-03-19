package model

import (
	"encoding/json"
	"strings"
)

type Version struct {
	Version int64 `yaml:"version"`
	Rules   Rule  `yaml:"rules"`
}

func (v *Version) Validate() bool {
	if (v.Version >= 0 && v.Version <= 7) && v.Rules.Validate() == true {
		return true
	}
	return false
}

func (v *Version) String() string {
	result, _ := json.Marshal(v)
	return string(result)
}
func (v *Version) GetBusicessEntity(business string) Table {
	for entityIdx := range v.Rules.Tables {
		if strings.ToLower(v.Rules.Tables[entityIdx].Name) == business {
			return v.Rules.Tables[entityIdx]
		}
	}
	return Table{}
}

func (v Version) GetTable(tableName string) Table {
	for tableIdx := range v.Rules.Tables {
		if v.Rules.Tables[tableIdx].Name == tableName {
			return v.Rules.Tables[tableIdx]
		}
	}
	return Table{}
}
