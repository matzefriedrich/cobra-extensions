package reflection

import "reflect"

type CommandReflector[T any] interface {
	ReflectCommandDescriptor(n T) CommandDescriptor
}

type valueItem struct {
	value     reflect.Value
	valueType reflect.Type
}
