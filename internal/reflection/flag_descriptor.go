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
	settingKey  string
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

// WithSettingKey sets the setting-key for the flag.
func (d FlagDescriptor) WithSettingKey(key string) FlagDescriptor {
	d.settingKey = key
	return d
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
		value, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		return d.SetValue(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		return d.SetValue(value)
	case reflect.Slice:
		return d.setSliceValueFromText(text)
	}
	return fmt.Errorf("unsupported flag type: %v", d.kind)
}

func (d *FlagDescriptor) setSliceValueFromText(text string) error {
	parts := strings.Split(text, ",")
	switch d.elementKind {
	case reflect.String:
		return d.SetValue(parts)
	case reflect.Int:
		values, err := parseIntSlice(parts)
		if err != nil {
			return err
		}
		return d.SetValue(values)
	case reflect.Int64:
		values, err := parseInt64Slice(parts)
		if err != nil {
			return err
		}
		return d.SetValue(values)
	case reflect.Bool:
		values, err := parseBoolSlice(parts)
		if err != nil {
			return err
		}
		return d.SetValue(values)
	}
	return fmt.Errorf("unsupported flag type: %v", d.kind)
}

func parseIntSlice(parts []string) ([]int, error) {
	var values []int
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		value, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func parseInt64Slice(parts []string) ([]int64, error) {
	var values []int64
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		value, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func parseBoolSlice(parts []string) ([]bool, error) {
	var values []bool
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		value, err := strconv.ParseBool(trimmed)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func invalidValueError() error {
	return errors.New(ErrorInvalidValue)
}
