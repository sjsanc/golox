package environment

import (
	"fmt"

	"github.com/sjsanc/golox/internal/errors"
	"github.com/sjsanc/golox/internal/token"
)

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() *Environment {
	return &Environment{
		Values: make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	fmt.Printf("Defining %s as %v\n", name, value)
	e.Values[name] = value
}

func (e *Environment) Get(name *token.Token) (interface{}, error) {
	if value, ok := e.Values[name.Lexeme]; ok {
		return value, nil
	}
	return nil, errors.NewRuntimeErr(name, "Undefined variable: "+name.Lexeme)
}
