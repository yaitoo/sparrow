package model

import (
	"encoding/json"

	"github.com/yaitoo/sparrow/db/util"
)

type Tag struct {
	ID     int      `yaml:"id"`
	Names  []string `yaml:names`
	Amount int      `yaml:"amount"`
	Date   string   `yaml:"date,omitempty"`
	nodes  []string
}

func (t *Tag) GetIdString() string {
	return util.IntLeft(t.ID, 2, "0")
}

func (t *Tag) String() string {
	result, _ := json.Marshal(t)
	return string(result)
}

func (t *Tag) IsNullObject() bool {
	return len(t.Names) == 0
}

func (t *Tag) GetHashValue(sequnce int64) int64 {
	return sequnce % int64(t.Amount)
}

func (t *Tag) Validate() bool {
	if t.ID >= 0 && t.ID <= 7 && len(t.Names) > 0 {
		return true
	}
	return false
}
