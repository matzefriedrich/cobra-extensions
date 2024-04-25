package reflection

import (
	"errors"
	"reflect"
)

type FlagDescriptor struct {
	name  string
	kind  reflect.Kind
	value reflect.Value
	usage string
}

func (d *FlagDescriptor) AsString() string {
	return d.value.String()
}

func (d *FlagDescriptor) AsInt64() int64 {
	return d.value.Int()
}

func (d *FlagDescriptor) AsBool() bool {
	return d.value.Bool()
}

func NewFlagDescriptor(name string, usage string, t reflect.Kind, v reflect.Value) FlagDescriptor {
	return FlagDescriptor{
		name:  name,
		usage: usage,
		kind:  t,
		value: v,
	}
}

func (d *FlagDescriptor) SetValue(value interface{}) error {
	switch d.kind {
	case reflect.String:
		s, ok := value.(string)
		if ok {
			d.value.SetString(s)
			return nil
		}
		return invalidValueError()
	case reflect.Int, reflect.Int64:
		n, ok := value.(int64)
		if ok {
			d.value.SetInt(n)
			return nil
		}
		return invalidValueError()
	case reflect.Bool:
		b, ok := value.(bool)
		if ok {
			d.value.SetBool(b)
			return nil
		}
		return invalidValueError()
	}

	return errors.New("type not supported. flag must be of type string, int, or bool")
}

func invalidValueError() error {
	return errors.New("the specified value does not match the flag type")
}
