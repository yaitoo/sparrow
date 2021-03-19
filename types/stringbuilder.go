package types

import "strings"

import "fmt"

//StringBuilder  包裝strings.Builder提供便利的操作函式庫
type StringBuilder struct {
	strings.Builder
}

//AppendLine 增加行
func (sb *StringBuilder) AppendLine(args ...interface{}) *StringBuilder {
	fmt.Fprintln(sb, fmt.Sprint(args...)+"\n")

	return sb
}

//AppendLinef 增加行
func (sb *StringBuilder) AppendLinef(format string, args ...interface{}) *StringBuilder {
	fmt.Fprintln(sb, fmt.Sprintf(format, args...)+"\n")

	return sb
}
