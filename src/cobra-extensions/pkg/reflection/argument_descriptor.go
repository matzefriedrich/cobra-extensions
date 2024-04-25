package reflection

import "reflect"

type ArgumentDescriptor struct {
	argumentIndex int
	value         reflect.Value
}

func (a *ArgumentDescriptor) ArgumentIndex() int {
	return a.argumentIndex
}

func (a *ArgumentDescriptor) SetValue(value string) error {
	a.value.SetString(value)
	return nil
}

type ArgumentsDescriptor struct {
	MinimumArgs int
	Args        []ArgumentDescriptor
}
