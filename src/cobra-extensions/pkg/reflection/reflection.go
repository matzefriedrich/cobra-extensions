package reflection

import (
	"github.com/matzefriedrich/cobra-extensions/internal"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"reflect"
)

type commandReflector[T any] struct {
}

// NewCommandReflector Creates a new CommandReflector instance.
func NewCommandReflector[T any]() CommandReflector[T] {
	return &commandReflector[T]{}
}

// ReflectCommandDescriptor Reflects all metadata from a command handler and returns a new CommandDescriptor instance.
func (r *commandReflector[T]) ReflectCommandDescriptor(n T) CommandDescriptor {

	var flags = make([]FlagDescriptor, 0)
	arguments := NewArgumentsDescriptorWith()

	value := reflect.ValueOf(n)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	valueType := reflect.TypeOf(value.Interface())
	valueTypeName := valueType.Name()
	use := internal.ExtractCommandUse(valueTypeName)
	shortDescriptionText := ""
	longDescriptionText := ""

	stack := internal.MakeStack[valueItem]()
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
			if fieldType == reflect.TypeOf(abstractions.CommandName{}) {
				use = flagName
				shortDescriptionText = field.Tag.Get("short")
				longDescriptionText = field.Tag.Get("long")
				continue
			}

			fieldValue := next.value.Field(i)

			if tryReflectArgumentsDescriptor(fieldType, fieldValue, arguments) {
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

func tryReflectArgumentsDescriptor(fieldType reflect.Type, fieldValue reflect.Value, target ArgumentsDescriptor) bool {

	m := ReflectedMember{value: fieldValue, memberType: fieldType}

	hasCommandArgs := false

	m.EnumerateFields(func(index int, field ReflectedField) {
		switch field.typeKind() {
		case reflect.String:
			if hasCommandArgs {
				target.With(Args(ArgumentDescriptor{value: field.value, argumentIndex: index - 1}))
			}
		case reflect.Struct:
			if field.isType(abstractions.CommandArgs{}) {
				compatible, ok := field.getInterfaceValue().(abstractions.CommandArgs)
				if ok {
					target.With(MinimumArgs(compatible.MinimumArgs))
					hasCommandArgs = true
				}
			}
		}
	})

	/* if fieldType.Kind() == reflect.Struct {
		fieldTypeNumFields := fieldType.NumField()
		for i := 0; i < fieldTypeNumFields; i++ {
			structFieldType := fieldType.Field(i)
			if structFieldType.Type == reflect.TypeOf(abstractions.CommandArgs{}) {
				structFieldValue := fieldValue.Field(i)
				compatible, ok := structFieldValue.Interface().(abstractions.CommandArgs)
				if ok {
					target.With(MinimumArgs(compatible.MinimumArgs))
					for j := i + 1; j < fieldTypeNumFields; j++ {
						kind := fieldType.Field(j).Type.Kind()
						if kind == reflect.String {
							arg := ArgumentDescriptor{value: fieldValue.Field(j), argumentIndex: j - 1}
							target.With(Args(arg))
						}
					}
					return true
				}
			}
		}
	}

	return false */
	return hasCommandArgs
}
