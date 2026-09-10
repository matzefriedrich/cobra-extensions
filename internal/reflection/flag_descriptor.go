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
func NewFlagDescriptor(name string, shorthand string, usage string, kind reflect.Kind, elementKind reflect.Kind, value reflect.Value) FlagDescriptor {
	return FlagDescriptor{
		name:        name,
		shorthand:   shorthand,
		usage:       usage,
		kind:        kind,
		elementKind: elementKind,
		value:       value,
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
		return d.setStringValue(value)
	case reflect.Int, reflect.Int64:
		return d.setInt64Value(value)
	case reflect.Bool:
		return d.setBoolValue(value)
	case reflect.Slice:
		return d.setSliceValue(value)
	}

	return errors.New(ErrorFlagTypeNotSupported)
}

func (d *FlagDescriptor) setStringValue(value interface{}) error {
	text, ok := value.(string)
	if !ok {
		return invalidValueError()
	}
	d.value.SetString(text)
	return nil
}

func (d *FlagDescriptor) setInt64Value(value interface{}) error {
	number, ok := value.(int64)
	if !ok {
		return invalidValueError()
	}
	d.value.SetInt(number)
	return nil
}

func (d *FlagDescriptor) setBoolValue(value interface{}) error {
	boolean, ok := value.(bool)
	if !ok {
		return invalidValueError()
	}
	d.value.SetBool(boolean)
	return nil
}

func (d *FlagDescriptor) setSliceValue(value interface{}) error {
	reflectedValue := reflect.ValueOf(value)
	if reflectedValue.Kind() != reflect.Slice {
		return invalidValueError()
	}
	d.value.Set(reflectedValue)
	return nil
}

// SetValueFromText sets the flag's value from its string representation.
func (d *FlagDescriptor) SetValueFromText(text string) error {
	switch d.kind {
	case reflect.String:
		return d.SetValue(text)
	case reflect.Int, reflect.Int64:
		return d.setTextAsInt64(text)
	case reflect.Bool:
		return d.setTextAsBool(text)
	case reflect.Slice:
		return d.setSliceValueFromText(text)
	}
	return fmt.Errorf("unsupported flag type: %v", d.kind)
}

func (d *FlagDescriptor) setTextAsInt64(text string) error {
	value, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return err
	}
	return d.SetValue(value)
}

func (d *FlagDescriptor) setTextAsBool(text string) error {
	value, err := strconv.ParseBool(text)
	if err != nil {
		return err
	}
	return d.SetValue(value)
}

func (d *FlagDescriptor) setSliceValueFromText(text string) error {
	parts := strings.Split(text, ",")
	switch d.elementKind {
	case reflect.String:
		return d.SetValue(parts)
	case reflect.Int:
		return d.setIntsFromParts(parts)
	case reflect.Int64:
		return d.setInt64sFromParts(parts)
	case reflect.Bool:
		return d.setBoolsFromParts(parts)
	}
	return fmt.Errorf("unsupported flag type: %v", d.kind)
}

func (d *FlagDescriptor) setIntsFromParts(parts []string) error {
	values, err := parseIntSlice(parts)
	if err != nil {
		return err
	}
	return d.SetValue(values)
}

func (d *FlagDescriptor) setInt64sFromParts(parts []string) error {
	values, err := parseInt64Slice(parts)
	if err != nil {
		return err
	}
	return d.SetValue(values)
}

func (d *FlagDescriptor) setBoolsFromParts(parts []string) error {
	values, err := parseBoolSlice(parts)
	if err != nil {
		return err
	}
	return d.SetValue(values)
}

func parseIntSlice(parts []string) ([]int, error) {
	return parseSlice(parts, strconv.Atoi)
}

func parseInt64Slice(parts []string) ([]int64, error) {
	return parseSlice(parts, func(part string) (int64, error) {
		return strconv.ParseInt(part, 10, 64)
	})
}

func parseBoolSlice(parts []string) ([]bool, error) {
	return parseSlice(parts, strconv.ParseBool)
}

func parseSlice[T any](parts []string, parse func(string) (T, error)) ([]T, error) {
	var values []T
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		value, err := parse(trimmed)
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
