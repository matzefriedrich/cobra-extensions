package reflection

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	ErrorFlagTypeNotSupported = "type not supported. flag must be of type string, int, or bool"
	ErrorInvalidValue         = "the specified instanceValue does not match the flag type"
)

// FlagDescriptor holds metadata for a command flag, including its name, type, value, and usage description.
type FlagDescriptor struct {
	name        string
	shorthand   string
	kind        reflect.Kind
	elementKind reflect.Kind
	value       reflect.Value
	usage       string
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

// NewFlagDescriptor creates a new FlagDescriptor given the flag's name, shorthand, usage description, type, and initial value.
func NewFlagDescriptor(name string, shorthand string, usage string, t reflect.Kind, et reflect.Kind, v reflect.Value) FlagDescriptor {
	return FlagDescriptor{
		name:        name,
		shorthand:   shorthand,
		usage:       usage,
		kind:        t,
		elementKind: et,
		value:       v,
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
	case reflect.Slice:
		v := reflect.ValueOf(value)
		if v.Kind() == reflect.Slice {
			d.value.Set(v)
			return nil
		}
		return invalidValueError()
	}

	return errors.New(ErrorFlagTypeNotSupported)
}

// SetValueFromText sets the flag's value from its string representation.
func (d *FlagDescriptor) SetValueFromText(text string) error {
	switch d.kind {
	case reflect.String:
		return d.SetValue(text)
	case reflect.Int, reflect.Int64:
		val, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		return d.SetValue(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		return d.SetValue(val)
	case reflect.Slice:
		parts := strings.Split(text, ",")
		switch d.elementKind {
		case reflect.String:
			return d.SetValue(parts)
		case reflect.Int:
			intParts := make([]int, len(parts))
			for i, p := range parts {
				v, err := strconv.Atoi(strings.TrimSpace(p))
				if err != nil {
					return err
				}
				intParts[i] = v
			}
			return d.SetValue(intParts)
		case reflect.Int64:
			intParts := make([]int64, len(parts))
			for i, p := range parts {
				v, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
				if err != nil {
					return err
				}
				intParts[i] = v
			}
			return d.SetValue(intParts)
		case reflect.Bool:
			boolParts := make([]bool, len(parts))
			for i, p := range parts {
				v, err := strconv.ParseBool(strings.TrimSpace(p))
				if err != nil {
					return err
				}
				boolParts[i] = v
			}
			return d.SetValue(boolParts)
		}
	}
	return fmt.Errorf("unsupported flag type: %v", d.kind)
}

func invalidValueError() error {
	return errors.New(ErrorInvalidValue)
}
