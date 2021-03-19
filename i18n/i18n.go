package i18n

import (
	"context"
	"fmt"
)

//T transfer code for current context
func T(ctx context.Context, code, defaultText string) string {
	//TODO:
	return defaultText
}

//Tf transfer code with format for current context
func Tf(ctx context.Context, code string, defaultFormat string, args ...interface{}) string {
	//TODO:

	return fmt.Sprintf(defaultFormat, args...)
}
