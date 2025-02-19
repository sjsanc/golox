package interpreter

import (
	"fmt"

	"github.com/sjsanc/golox/internal/environment"
	"github.com/sjsanc/golox/internal/stmt"
)

type Function struct {
	declaration stmt.Function
	closure     *environment.Environment
}

func NewFunction(declaration stmt.Function, closure *environment.Environment) *Function {
	return &Function{
		declaration: declaration,
		closure:     closure,
	}
}

func (f *Function) Arity() int {
	return len(f.declaration.Params)
}

func (f *Function) Call(interpreter *Interpreter, args []interface{}) (interface{}, error) {
	env := environment.NewEnvironment(f.closure)
	for i, param := range f.declaration.Params {
		env.Define(param.Lexeme, args[i])
	}

	val, err := interpreter.executeBlock(f.declaration.Body, env)
	if err != nil {
		return nil, err
	}

	return val.Value, nil
}

func (f *Function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}
