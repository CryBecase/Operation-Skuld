package env

import "fmt"

const (
	Local      Env = "local"
	Dev        Env = "dev"
	Test       Env = "test"
	Staging    Env = "staging"
	Production Env = "production"
)

var valid = []Env{Local, Dev, Test, Staging, Production}

type Env string

func New(s string) Env {
	e := Env(s)
	if !e.Valid() {
		panic(fmt.Sprintf("env must one of %v, but get %s", valid, s))
	}

	return e
}

func (e Env) IsLocal() bool {
	return e == Local
}

func (e Env) IsDev() bool {
	return e == Dev
}

func (e Env) IsTest() bool {
	return e == Test
}

func (e Env) IsStaging() bool {
	return e == Staging
}

func (e Env) IsProduction() bool {
	return e == Production
}

func (e Env) Valid() bool {
	for _, v := range valid {
		if e == v {
			return true
		}
	}
	return false
}
