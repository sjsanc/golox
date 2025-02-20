package tw

type Class struct {
	name string
}

func NewClass(name string) *Class {
	return &Class{
		name: name,
	}
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Arit() int {
	return 0
}

func (c *Class) Call(interpreter *Interpreter, args []interface{}) (interface{}, error) {
	return NewInstance(c), nil
}
