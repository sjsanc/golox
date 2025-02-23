package tw

import "fmt"

type Instance struct {
	class  *Class
	fields map[string]interface{}
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class:  class,
		fields: make(map[string]interface{}),
	}
}

func (i *Instance) Get(name *Token) (interface{}, error) {
	if value, ok := i.fields[name.lexeme]; ok {
		return value, nil
	}
	method := i.class.FindMethod(name.lexeme)
	if method != nil {
		return method.Bind(i), nil
	}
	return nil, fmt.Errorf("undefined property '%s'", name.lexeme)
}

func (i *Instance) Set(name *Token, value interface{}) {
	i.fields[name.lexeme] = value
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}
