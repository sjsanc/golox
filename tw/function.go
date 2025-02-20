package tw

import (
	"fmt"
)

type FunctionType string

const (
	FunctionNone     FunctionType = "none"
	FunctionFunction FunctionType = "function"
)

type Function struct {
	declaration FunctionStmt
	closure     *Environment
}

func NewFunction(declaration FunctionStmt, closure *Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.params)
}

func (f *Function) Call(interpreter *Interpreter, args []interface{}) (interface{}, error) {
	env := NewEnvironment(f.closure)
	for i, param := range f.declaration.params {
		env.Define(param.lexeme, args[i])
	}

	val, err := interpreter.executeBlock(f.declaration.body, env)
	if err != nil {
		return nil, err
	}

	return val.value, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}
