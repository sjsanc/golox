package tw

type Class struct {
	name    string
	methods map[string]*Function
}

func NewClass(name string) *Class {
	return &Class{
		name:    name,
		methods: make(map[string]*Function),
	}
}

func (c *Class) FindMethod(name string) *Function {
	if method, ok := c.methods[name]; ok {
		return method
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
