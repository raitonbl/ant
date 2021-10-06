package internal

type ProjectContext interface {
	GetDescriptor() string
}

type DefaultContext struct {
	Descriptor string
}

func (instance *DefaultContext) GetDescriptor() string {
	return instance.Descriptor
}
