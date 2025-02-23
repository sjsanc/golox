package tw

import (
	"fmt"
)

type Function struct {
	declaration   *FunctionStmt
	closure       *Environment
	isInitializer bool
}

func NewFunction(declaration *FunctionStmt, closure *Environment, isInitializer bool) *Function {
	return &Function{
		declaration:   declaration,
		closure:       closure,
		isInitializer: isInitializer,
	}
}

func (f *Function) Bind(instance *Instance) *Function {
	env := NewEnvironment(f.closure)
	env.Define("this", instance)
	return NewFunction(f.declaration, env, f.isInitializer)
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
		if f.isInitializer {
			return f.closure.GetAt(0, "this")
		}
		return nil, err
	}

	if f.isInitializer {
		return f.closure.GetAt(0, "this")
	}

	return val.value, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.name.lexeme)
}
