package types

import (
	"net/url"
	"sort"
	"strings"
)

//Field 表单字段
type Field struct {
	value string
	quote bool
}

//Form 表单
type Form struct {
	Keys   []string
	Fields map[string]*Field
}

//NewForm 创建Form对象
func NewForm() Form {
	sm := Form{}
	sm.Keys = make([]string, 0, 20)
	sm.Fields = make(map[string]*Field)

	return sm
}

//AddJSON 加入JSON字段数据
func (m *Form) AddJSON(key, value string, quote bool) *Form {
	m.Keys = append(m.Keys, key)
	m.Fields[key] = &Field{value: value, quote: quote}

	return m
}

//Add 添加字段数据
func (m *Form) Add(key, value string) *Form {
	m.Keys = append(m.Keys, key)
	m.Fields[key] = &Field{value: value, quote: true}

	return m
}

//Update 更新字段值
func (m *Form) Update(key, value string, quote bool) *Form {

	m.Fields[key] = &Field{value: value, quote: quote}

	return m
}

//Get 获取指定名称字段值
func (m *Form) Get(key string) string {
	val, ok := m.Fields[key]
	if ok {
		return val.value
	}

	return ""
}

//SortAndJoin 按字段名排序并且按指定字符做拼接“字段名称=字段值”，空值忽略
func (m *Form) SortAndJoin(sep string) string {

	sort.Strings(m.Keys)
	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		if len(value.value) > 0 {
			list = append(list, key+"="+value.value)
		}
	}

	return strings.Join(list, sep)
}

//SortAndJoinAll 按字段名排序并且按指定字符做拼接“字段名称=字段值”，空值不忽略
func (m *Form) SortAndJoinAll(sep string) string {

	sort.Strings(m.Keys)
	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		list = append(list, key+"="+value.value)
	}

	return strings.Join(list, sep)
}

//Join 按字段加入顺序，按指定字符做拼接“字段名称=字段值”，空值忽略
func (m *Form) Join(sep string) string {

	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		if len(value.value) > 0 {
			list = append(list, key+"="+value.value)
		}
	}

	return strings.Join(list, sep)
}

//JoinAll 按字段加入顺序，按指定字符做拼接“字段名称=字段值”，空值不忽略
func (m *Form) JoinAll(sep string) string {

	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		list = append(list, key+"="+value.value)
	}

	return strings.Join(list, sep)
}

//ToJSON  生成JSON数字串
func (m *Form) ToJSON() string {

	sort.Strings(m.Keys)
	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		if value.value != "" {
			if value.quote {
				list = append(list, "\""+key+"\":\""+strings.Replace(value.value, "\"", "\\\"", -1)+"\"")
			} else {
				list = append(list, "\""+key+"\":"+value.value)
			}

		}
	}

	return "{" + strings.Join(list, ",") + "}"
}

//ToForm 转成原始map对象
func (m *Form) ToForm() map[string]string {

	list := make(map[string]string)
	for k, v := range m.Fields {
		list[k] = v.value
	}

	return list
}

//EncodeURL 转成URL编码格式
func (m *Form) EncodeURL() string {

	values := &url.Values{}

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		if len(value.value) > 0 {
			values.Add(key, value.value)
		}
	}

	return values.Encode()
}

//JoinValues 按字段加入顺序，按指定字符做拼接字段值，空值忽略
func (m *Form) JoinValues(sep string) string {

	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		list = append(list, value.value)
	}

	return strings.Join(list, sep)
}

//JoinFunc 自定义拼接
func (m *Form) JoinFunc(each func(key, value string) (string, bool), sep string, sortFirst bool) string {
	if sortFirst {
		sort.Strings(m.Keys)
	}

	list := make([]string, 0, len(m.Keys))

	for _, key := range m.Keys {
		value, _ := m.Fields[key]
		v, ok := each(key, value.value)
		if ok {
			list = append(list, v)
		}

	}

	return strings.Join(list, sep)
}
