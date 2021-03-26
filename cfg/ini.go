// Copyright (c) 2021 Yaitoo.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cfg

import (
	"bufio"
	"context"
	"log"
	"strconv"
	"strings"
)

//Inifile represents all data from an INI file. see https://en.wikipedia.org/wiki/INI_file
// NOTE:
//   + Global properties will be ignored.
//   + Section and property names are not case sensitive. Both of them are converted to lower-case. And space is trimed in leading and trailing.
//   + The value of bool key is not case sensitive. "True" is same as "true". Space is trimed too.
type Inifile struct {
	sections map[string]*Section
}

//TryParse try parse string. if it is ok, Infile will be updated.
func (i *Inifile) TryParse(ctx context.Context, data string) {

	sections := make(map[string]*Section)

	var section *Section

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		if strings.HasPrefix(text, "#") || strings.HasPrefix(text, ";") {
			continue
		}

		if strings.HasPrefix(text, "[") {
			sectionName := strings.ToLower(strings.TrimSpace(text[1 : len(text)-1]))
			if section == nil {
				section = &Section{}
				section.Name = sectionName
				section.values = make(map[string]string)
				sections[section.Name] = section

			} else if section.Name != sectionName { //go to next section
				section = &Section{}
				section.Name = sectionName
				section.values = make(map[string]string)
				sections[sectionName] = section
			}

			continue
		}

		i := strings.Index(text, "=")
		if i > 0 {
			name := strings.ToLower(strings.TrimSpace(text[:i]))
			value := text[i+1:]

			if section == nil {
				log.Printf("[inifile]section is missing: %s\n", text)
				continue
			}

			section.values[name] = value
		} else {
			log.Printf("[inifile]'=' is missing: %s\n", text)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("[inifile]invalid ini data: %s\n", data)
	}

	i.sections = sections

}

//Section get section with name.
func (i *Inifile) Section(name string) *Section {
	if i == nil {
		return nil
	}

	if i.sections == nil {
		return nil
	}

	s, ok := i.sections[strings.ToLower(name)]
	if ok {
		return s
	}

	return nil
}

//Section  information associated to a section in a INI File
type Section struct {
	//Name the name of section
	Name string

	values map[string]string
}

//Keys return key list
func (s *Section) Keys() []string {
	n := len(s.values)
	keys := make([]string, n)

	i := 0
	for k := range s.values {
		keys[i] = k
		i++
	}

	return keys
}

//Value get string with key. return defaultValue if it doesn't exists
func (s *Section) Value(key string, defaultValue string) string {
	if s == nil || s.values == nil {
		return defaultValue
	}

	v, ok := s.values[strings.ToLower(strings.TrimSpace(key))]
	if ok {
		return v
	}

	return defaultValue
}

//ValueInt get int with key. return defaultValue if it doesn't exists or is invalid int
func (s *Section) ValueInt(key string, defaultValue int) int {
	v := strings.TrimSpace(s.Value(key, ""))

	if v == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 10, 0)
	if err != nil {
		log.Printf("Unable to interpret [%s]%s=%s as a int\n", s.Name, key, v)
		return defaultValue
	}

	return int(i)
}

//ValueInt32 get int32 with key. return defaultValue if it doesn't exists or is invalid int64
func (s *Section) ValueInt32(key string, defaultValue int32) int32 {
	v := strings.TrimSpace(s.Value(key, ""))

	if v == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		log.Printf("Unable to interpret [%s]%s=%s as a int64\n", s.Name, key, v)
		return defaultValue
	}

	return int32(i)
}

//ValueInt64 get int64 with key. return defaultValue if it doesn't exists or is invalid int64
func (s *Section) ValueInt64(key string, defaultValue int64) int64 {
	v := strings.TrimSpace(s.Value(key, ""))

	if v == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Printf("Unable to interpret [%s]%s=%s as a int64\n", s.Name, key, v)
		return defaultValue
	}

	return int64(i)
}

//ValueFloat32 get float32 with key. return defaultValue if it doesn't exists or is invalid float32
func (s *Section) ValueFloat32(key string, defaultValue float32) float32 {
	v := strings.TrimSpace(s.Value(key, ""))

	if v == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(v, 32)
	if err != nil {
		log.Printf("Unable to interpret [%s]%s=%s as a float32\n", s.Name, key, v)
		return defaultValue
	}

	return float32(i)
}

//ValueFloat64 get float64 with key. return defaultValue if it doesn't exists or is invalid float32
func (s *Section) ValueFloat64(key string, defaultValue float64) float64 {
	v := strings.TrimSpace(s.Value(key, ""))

	if v == "" {
		return defaultValue
	}

	i, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Printf("Unable to interpret [%s]%s=%s as a float64\n", s.Name, key, v)
		return defaultValue
	}

	return float64(i)
}

//ValueBool get bool with key. return defaultValue if it doesn't exists or is invalid value. valid values: 0/1, on/off, true/false.
func (s *Section) ValueBool(key string, defaultValue bool) bool {
	v := strings.ToLower(strings.TrimSpace(s.Value(key, "")))

	if v == "1" || v == "on" || v == "true" {
		return true
	}

	if v == "0" || v == "off" || v == "false" {
		return false
	}

	return defaultValue

}
