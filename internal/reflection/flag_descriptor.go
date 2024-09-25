package reflection

import (
	"errors"
	"reflect"
)

const (
	ErrorFlagTypeNotSupported = "type not supported. flag must be of type string, int, or bool"
	ErrorInvalidValue         = "the specified instanceValue does not match the flag type"
)

// FlagDescriptor holds metadata for a command flag, including its name, type, value, and usage description.
type FlagDescriptor struct {
	name  string
	kind  reflect.Kind
	value reflect.Value
	usage string
}

// AsString returns the string representation of the flag's value.
func (d *FlagDescriptor) AsString() string {
	return d.value.String()
}

// AsInt64 returns the int64 representation of the flag's value.
func (d *FlagDescriptor) AsInt64() int64 {
	return d.value.Int()
}

// AsBool returns the boolean representation of the flag's value.
func (d *FlagDescriptor) AsBool() bool {
	return d.value.Bool()
}

// NewFlagDescriptor creates a new FlagDescriptor given the flag's name, usage description, type, and initial value.
func NewFlagDescriptor(name string, usage string, t reflect.Kind, v reflect.Value) FlagDescriptor {
	return FlagDescriptor{
		name:  name,
		usage: usage,
		kind:  t,
		value: v,
	}
}

// SetValue sets the value of a flag based on its type (string, int64, or bool) and returns an error if the type is unsupported.
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

	return errors.New(ErrorFlagTypeNotSupported)
}

func invalidValueError() error {
	return errors.New(ErrorInvalidValue)
}
