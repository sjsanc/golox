package tw

type Class struct {
	name       string
	superclass *Class
	methods    map[string]*Function
}

func NewClass(name string, superclass *Class, methods map[string]*Function) *Class {
	return &Class{
		name:       name,
		superclass: superclass,
		methods:    methods,
	}
}

func (c *Class) FindMethod(name string) *Function {
	if method, ok := c.methods[name]; ok {
		return method
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil
}

func (c *Class) String() string {
	return c.name
}

func (c *Class) Arity() int {
	initializer := c.FindMethod("init")
	if initializer != nil {
		return initializer.Arity()
	}
	return 0
}

func (c *Class) Call(interpreter *Interpreter, args []interface{}) (interface{}, error) {
	instance := NewInstance(c)
	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(interpreter, args)
	}
	return instance, nil
}
