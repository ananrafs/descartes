package common

type TypeChecker struct {
	Type string `json:"type"`
}

type TypeCheckerItf interface {
	GetType() string
}
