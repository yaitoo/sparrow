package model

import (
	"encoding/json"
	"time"

	"github.com/yaitoo/sparrow/db/util"
)

type Table struct {
	Name     string `yaml:name`
	Tags     []Tag  `yaml:tags`
	Key      string `yaml:key`
	TimeZone string `yaml:"timeZone"`
}

func (t *Table) Validate() bool {
	if t.Name != "" && t.Key != "" {
		for i := 0; i < len(t.Tags); i++ {
			if t.Tags[i].Validate() == false {
				return false
			}
		}
		return true
	}
	return false
}

/* func (e *Table) Validating() bool {
	mapIdHash := make(map[int]int)
	for _, v := range e.Tags {
		if val, ok := mapIdHash[v.ID]; ok {
			if val != v.Amount {
				return false
			}
		} else {
			mapIdHash[v.ID] = v.Amount
		}
	}
	return true
}
*/
func (e *Table) String() string {
	result, _ := json.Marshal(e)
	return string(result)
}

func (e *Table) IsEmptyTimeZone() bool {
	return e.TimeZone == ""
}

func (e *Table) IsNullObject() bool {
	return e.Name == ""
}

func (e *Table) GetTable(bizName string) Tag {
	for idx := range e.Tags {
		for tagIdx := range e.Tags[idx].Names {
			if e.Tags[idx].Names[tagIdx] == bizName {
				return e.Tags[idx]
			}
		}
	}
	return Tag{}
}

func (t Table) GetTag(tagId int64) Tag {
	for idx := range t.Tags {
		if int64(t.Tags[idx].ID) == tagId {
			return t.Tags[idx]
		}
	}
	return Tag{}
}

func (e *Table) GetTime(targetTIme time.Time) (util.TimeUtil, error) {
	return util.ConvertWIthTimeZone(targetTIme, e.TimeZone)
}
