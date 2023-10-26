package pkg

import (
	"github.com/matzefriedrich/cobra-extensions/internal"
	"reflect"

	"github.com/spf13/cobra"
)

type valueItem struct {
	value     reflect.Value
	valueType reflect.Type
}

func ReflectCommandDescriptor[T any](n T) CommandDescriptor {

	var flags = make([]FlagDescriptor, 0)

	value := reflect.ValueOf(n)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	valueType := reflect.TypeOf(value.Interface())
	valueTypeName := valueType.Name()
	use := internal.ExtractCommandUse(valueTypeName)

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
			if fieldType == reflect.TypeOf(CommandName{}) {
				use = flagName
				continue
			}

			fieldValue := next.value.Field(i)

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

	return NewCommandDescriptor(use, flags)
}

func UnmarshalCommand(source *cobra.Command, desc CommandDescriptor) {
	for _, f := range desc.flags {
		flags := source.Flags()
		flagName := f.name
		switch f.kind {
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
}
