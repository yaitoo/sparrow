package db

type Validation interface {
	IsValidate() bool
}

type ValiadtionCallback interface {
	Validate(func() bool) bool
}
