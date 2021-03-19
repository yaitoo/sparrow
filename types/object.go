package types

import "context"

//Object base class for any object
type Object interface {
	IsNull() bool
}

//Contexter  an object which accepts context
type Contexter interface {
	SetContext(ctx context.Context)
}
