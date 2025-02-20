package tw

type Instance struct {
	class *Class
}

func NewInstance(class *Class) *Instance {
	return &Instance{
		class: class,
	}
}

func (i *Instance) String() string {
	return i.class.name + " instance"
}
