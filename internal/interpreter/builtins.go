package interpreter

import "time"

type ClockBuiltin struct {
	arity int
}

func (c *ClockBuiltin) Arity() int {
	return c.arity
}

func (c *ClockBuiltin) Call(interpreter *Interpreter, args []interface{}) interface{} {
	return time.Now().Unix()
}

func (c *ClockBuiltin) String() string {
	return "<native fn>"
}
