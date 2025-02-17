package environment

import (
	"github.com/sjsanc/golox/internal/errors"
	"github.com/sjsanc/golox/internal/token"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewGlobalEnvironment() *Environment {
	return &Environment{
		Values: make(map[string]interface{}),
	}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Get(name *token.Token) (interface{}, error) {
	if value, ok := e.Values[name.Lexeme]; ok {
		return value, nil
	}
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}
	return nil, errors.NewRuntimeErr(name, "Undefined variable: "+name.Lexeme)
}

func (e *Environment) Assign(name *token.Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}
	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}
	return errors.NewRuntimeErr(name, "Undefined variable: "+name.Lexeme)
}
