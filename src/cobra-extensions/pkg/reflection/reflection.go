package reflection

import (
	"github.com/matzefriedrich/cobra-extensions/internal"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"reflect"

	"github.com/spf13/cobra"
)

type commandReflector[T any] struct {
}

// NewCommandReflector Creates a new CommandReflector instance.
func NewCommandReflector[T any]() CommandReflector[T] {
	return &commandReflector[T]{}
}

func (r *commandReflector[T]) ReflectCommandDescriptor(n T) CommandDescriptor {

	var flags = make([]FlagDescriptor, 0)
	var arguments = ArgumentsDescriptor{
		Args: make([]ArgumentDescriptor, 0),
	}

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

			if tryReflectArgumentsDescriptor(fieldType, fieldValue, &arguments) {
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

func tryReflectArgumentsDescriptor(fieldType reflect.Type, fieldValue reflect.Value, target *ArgumentsDescriptor) bool {

	if fieldType.Kind() == reflect.Struct {
		fieldTypeNumFields := fieldType.NumField()
		for i := 0; i < fieldTypeNumFields; i++ {
			structFieldType := fieldType.Field(i)
			if structFieldType.Type == reflect.TypeOf(abstractions.CommandArgs{}) {
				structFieldValue := fieldValue.Field(i)
				compatible, ok := structFieldValue.Interface().(abstractions.CommandArgs)
				if ok {
					target.MinimumArgs = compatible.MinimumArgs
					for j := i + 1; j < fieldTypeNumFields; j++ {
						kind := fieldType.Field(j).Type.Kind()
						if kind == reflect.String {
							arg := ArgumentDescriptor{value: fieldValue.Field(j), argumentIndex: j - 1}
							target.Args = append(target.Args, arg)
						}
					}
					return true
				}
			}
		}
	}

	return false
}

func UnmarshalCommand(source *cobra.Command, desc CommandDescriptor, args ...string) {

	for _, f := range desc.Flags() {
		flags := source.Flags()
		flagName := f.Name()
		switch f.Kind() {
		case reflect.String:
			value, _ := flags.GetString(flagName)
			_ = f.SetValue(value)
		case reflect.Int, reflect.Int64:
			value, _ := flags.GetInt64(flagName)
			_ = f.SetValue(value)
		case reflect.Bool:
			value, _ := flags.GetBool(flagName)
			_ = f.SetValue(value)
		}
	}

	arguments := desc.Arguments()
	for _, a := range arguments.Args {
		index := a.argumentIndex
		if index < len(args) {
			value := args[index]
			_ = a.SetValue(value)
		}
	}
}
