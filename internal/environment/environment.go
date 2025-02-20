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

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	env := e.ancestor(distance)
	if value, ok := env.Values[name]; ok {
		return value, nil
	}
	return nil, errors.NewRuntimeErr(nil, "Undefined variable: "+name)
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

func (e *Environment) AssignAt(distance int, name *token.Token, value interface{}) error {
	env := e.ancestor(distance)
	env.Values[name.Lexeme] = value
	return nil
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Enclosing
	}
	return env
}
