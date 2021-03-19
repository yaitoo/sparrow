package models

//Field 验证规则
type Field struct {
	//Rule 规则类型
	Rule string
	//Message i18n提示信息
	Message map[string]string
}

//Form 验证表单，key=栏位名称，值=验证规则集合
type Form map[string][]Field
