package reflection

import (
	"github.com/matzefriedrich/cobra-extensions/internal/utils"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"reflect"
)

type commandReflector[T any] struct {
}

// NewCommandReflector Creates a new CommandReflector instance.
func NewCommandReflector[T any]() types.CommandReflector[T] {
	return &commandReflector[T]{}
}

// ReflectCommandDescriptor Reflects all metadata from a command handler and returns a new CommandDescriptor instance.
func (r *commandReflector[T]) ReflectCommandDescriptor(n T) types.CommandDescriptor {

	var flags = make([]FlagDescriptor, 0)
	arguments := NewArgumentsDescriptorWith()

	value := reflect.ValueOf(n)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	valueType := reflect.TypeOf(value.Interface())
	valueTypeName := valueType.Name()
	use := utils.ExtractCommandUse(valueTypeName)
	shortDescriptionText := ""
	longDescriptionText := ""

	stack := utils.MakeStack[valueItem]()
	stack.Push(valueItem{value: value, valueType: valueType})

	for {

		if stack.IsEmpty() {
			break
		}

		next := stack.Pop()

		numFields := next.value.NumField()
		for i := 0; i < numFields; i++ {

			field := next.valueType.Field(i)
			isExportedField := field.PkgPath == ""

			flagName := field.Tag.Get("flag")

			fieldType := field.Type
			if fieldType == reflect.TypeOf(types.CommandName{}) {
				use = flagName
				shortDescriptionText = field.Tag.Get("short")
				longDescriptionText = field.Tag.Get("long")
				continue
			}

			fieldValue := next.value.Field(i)

			m := ReflectedObject{instanceValue: fieldValue, objectType: fieldType}
			if tryReflectArgumentsDescriptor(m, arguments) {
				continue
			}

			isEmbeddedField := field.Anonymous
			if isEmbeddedField {
				embeddedValue := fieldValue
				embeddedType := fieldType
				stack.Push(valueItem{value: embeddedValue, valueType: embeddedType})
				continue
			}

			if isExportedField {
				usage := field.Tag.Get("usage")
				fieldTypeKind := fieldType.Kind()
				desc := NewFlagDescriptor(flagName, usage, fieldTypeKind, fieldValue)
				flags = append(flags, desc)
			}
		}
	}

	return NewCommandDescriptor(use, shortDescriptionText, longDescriptionText, flags, arguments)
}

func tryReflectArgumentsDescriptor(m ReflectedObject, target types.ArgumentsDescriptor) bool {

	hasCommandArgs := false

	m.EnumerateFields(func(index int, field ReflectedField) {
		fieldTypeKind := field.typeKind()
		switch fieldTypeKind {
		case reflect.String:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Bool:
			if hasCommandArgs {
				descriptor := ArgumentDescriptor{typeKind: fieldTypeKind, value: field.value, argumentIndex: index - 1}
				target.With(Args(descriptor))
			}
		case reflect.Interface:
		case reflect.Struct:
			if field.isType(types.CommandArgs{}) {
				compatible, ok := field.getInterfaceValue().(types.CommandArgs)
				if ok {
					target.With(MinimumArgs(compatible.MinimumArgs))
					hasCommandArgs = true
				}
			}
		}
	})

	return hasCommandArgs
}
