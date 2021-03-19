package linq

import (
	"strings"

	"github.com/yaitoo/sparrow/types"
)

type Strings struct {
	items []string
}

type StringPredicate func(s string) bool
type StringSelect func(s string) string

func FromStrings(list []string) Strings {
	return Strings{items: list}
}

func FromString(s, sep string) Strings {
	if types.IsEmpty(s) {
		return Strings{}
	}

	return Strings{items: strings.Split(s, sep)}
}

func (sl Strings) Where(sp StringPredicate) Strings {
	if len(sl.items) > 0 {
		items := make([]string, 0, len(sl.items))

		for _, item := range sl.items {
			if sp(item) {
				items = append(items, item)
			}
		}

		sl.items = items

		return sl
	}
	return sl
}

func (sl Strings) Select(ss StringSelect) Strings {
	if len(sl.items) > 0 {
		items := make([]string, 0, len(sl.items))

		for _, item := range sl.items {
			items = append(items, ss(item))
		}

		sl.items = items

		return sl
	}
	return sl
}

func (sl Strings) Count() int {
	return len(sl.items)
}

func (sl Strings) ToArray() []string {
	if sl.items != nil {

		return sl.items
	}
	return []string{}
}
