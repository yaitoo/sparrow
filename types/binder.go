package types

import "context"

//Binder 物件綁定器
type Binder interface {
	//Bind 自動綁定數據到對應到欄位上
	Bind(ctx context.Context, data map[string]string) error
}
