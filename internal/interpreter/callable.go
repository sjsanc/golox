package interpreter

type Callable interface {
	Arity() int
	Call(interpreter *Interpreter, args []interface{}) (interface{}, error)
}
