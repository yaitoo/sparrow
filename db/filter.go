package db

import (
	"github.com/yaitoo/sparrow/types"
)

//PagingIndex 分頁頁碼默認變數名稱
const PagingIndex = "_pi"

//PagingSize 分頁大小默認變數名稱
const PagingSize = "_ps"

//Filter 過濾器，自動類型轉換，自動索引優先
type Filter struct {
	fields     map[string]string
	converters map[string]func(s string) interface{}
	vars       map[string]interface{}
}

//NewFilter 創建過濾器
func NewFilter(values map[string]string) *Filter {
	f := &Filter{
		fields:     values,
		converters: make(map[string]func(s string) interface{}),
	}

	return f
}

//ToTake 提取分頁欄位
func (f *Filter) ToTake() (int, int) {
	if f == nil || f.fields == nil {
		return 0, 0
	}

	pi := f.fields[PagingIndex]
	ps := f.fields[PagingSize]

	return types.Atoi(pi, 0), types.Atoi(ps, 0)

}

//ToVars 型別轉換
func (f *Filter) ToVars() map[string]interface{} {
	if f == nil {
		return nil
	}

	if f.vars != nil {
		return f.vars
	}

	vars := make(map[string]interface{})

	for k, v := range f.fields {
		if types.IsNotEmpty(v) {
			if f.converters != nil {
				c, ok := f.converters[k]
				if ok {
					vars["a"] = c(v)
					continue
				}
			}

			vars[k] = v
		}

	}

	f.vars = vars

	return f.vars
}

//Map 註冊型別轉換
func (f *Filter) Map(converter func(s string) interface{}, fields ...string) *Filter {
	if f == nil {
		return f
	}

	if f.converters == nil {
		f.converters = make(map[string]func(s string) interface{})
	}

	for _, field := range fields {
		f.converters[field] = converter
	}

	return f
}
