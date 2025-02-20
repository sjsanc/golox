package tw

import "fmt"

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewGlobalEnvironment() *Environment {
	return &Environment{
		values: make(map[string]interface{}),
	}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name *Token) (interface{}, error) {
	if value, ok := e.values[name.lexeme]; ok {
		return value, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("undefined variable: %s", name)
}

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	env := e.ancestor(distance)
	if value, ok := env.values[name]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("undefined variable: %s", name)
}

func (e *Environment) Assign(name *Token, value interface{}) error {
	if _, ok := e.values[name.lexeme]; ok {
		e.values[name.lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return fmt.Errorf("undefined variable: %s", name)
}

func (e *Environment) AssignAt(distance int, name *Token, value interface{}) error {
	env := e.ancestor(distance)
	env.values[name.lexeme] = value
	return nil
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}
