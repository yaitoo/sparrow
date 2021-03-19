package model

import (
	"encoding/json"
	"strings"
)

type Config struct {
	maxVer   Version
	Database Database  `yaml:"database"`
	Versions []Version `yaml:"versions"`
}

func (c *Config) Validate() bool {
	if c.Database.Validate() == true {
		for i := 0; i < len(c.Versions); i++ {
			if c.Versions[i].Validate() == false {
				return false
			}
		}
		return true
	}
	return false
}

func (c *Config) String() string {
	result, _ := json.Marshal(c)
	return string(result)
}

func (c *Config) GetLastVersionEntity(entityName string) Table {
	for i := len(c.Versions) - 1; i >= 0; i-- {
		return c.Versions[i].GetBusicessEntity(entityName)
	}
	return Table{}
}

func (c *Config) GetNewestVersion() Version {
	for _, version := range c.Versions {
		if version.Version >= c.maxVer.Version {
			c.maxVer = version
		}
	}
	return c.maxVer
}

func (c *Config) GetSpecificVersion(version int64) Version {
	for idx := range c.Versions {
		if c.Versions[idx].Version == version {
			return c.Versions[idx]
		}
	}
	return Version{}
}

func (vc *Version) GetVer() int64 {
	return vc.Version
}

func (vc Version) GetBusicessRule(tableName string) Rule {
	var result Rule = Rule{
		//	NodesInfo: NewNodeInfo(),
	}

	for entityIdx := range vc.Rules.Tables {
		if strings.ToLower(vc.Rules.Tables[entityIdx].Name) == tableName {
			result.Tables = []Table{
				vc.Rules.Tables[entityIdx],
			}
			result.Databases = vc.Rules.Databases
			return result
		}
	}

	return Rule{}
}

func (c *Config) GetDefaultDatabse() Database {
	return c.Database
}

func createEntityNode(ruleIdx int, nodeName string, tableID int, vr *Version) {
	for entitiyIdx := range vr.Rules.Tables {
		for tableIdx := range vr.Rules.Tables[entitiyIdx].Tags {
			if vr.Rules.Tables[entitiyIdx].Tags[tableIdx].ID == tableID {
				vr.Rules.Tables[entitiyIdx].Tags[tableIdx].nodes =
					append(vr.Rules.Tables[entitiyIdx].Tags[tableIdx].nodes, nodeName)
			}
		}
	}
}

func (c *Config) GetDatabase(version, dbID int64, tableName string) Database {
	rule := c.GetSpecificVersion(version).GetBusicessRule(tableName)
	if rule.IsEmptyObject() == true {
		return c.GetDefaultDatabse()
	}
	return rule.GetDatabase(dbID)
}
