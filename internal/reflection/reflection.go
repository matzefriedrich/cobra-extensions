package reflection

import (
	"errors"
	"reflect"

	"github.com/matzefriedrich/cobra-extensions/internal/utils"
	"github.com/matzefriedrich/cobra-extensions/internal/utils/ternary"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
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
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	valueType := reflect.TypeOf(value.Interface())
	valueTypeName := valueType.Name()
	use := utils.ExtractCommandUse(valueTypeName)
	shortHelpText := ""
	longHelpText := ""

	stack := utils.MakeStack[valueItem]()
	stack.Push(valueItem{value: value, valueType: valueType})

	for !stack.IsEmpty() {

		next := stack.Pop()

		numFields := next.value.NumField()
		for i := range numFields {

			field := next.valueType.Field(i)
			isExportedField := field.PkgPath == ""

			fieldType := field.Type
			//nolint:staticcheck // required for legacy functionality
			if fieldType == reflect.TypeFor[types.CommandName]() || reflect.TypeFor[types.BaseCommand]() == fieldType {
				tag, tagErr := reflectCobraXCommand(field)
				if tagErr != nil && errors.Is(tagErr, ErrCobraXCommandNotFound) {
					tag, tagErr = reflectLegacyCommand(field)
					if tagErr != nil && errors.Is(tagErr, ErrCobraXLegacyTagsNotFound) {
						continue
					}
				}
				if tag != nil {
					use = ternary.ValueOrDefault(tag.Use, ternary.NotNilOrWhitespace, use)
					shortHelpText = ternary.ValueOrDefault(tag.Help, ternary.NotNilOrWhitespace, shortHelpText)
					longHelpText = ternary.ValueOrDefault(tag.Description, ternary.NotNilOrWhitespace, longHelpText)
				}
				//nolint:staticcheck // required for legacy functionality
				if fieldType == reflect.TypeFor[types.CommandName]() {
					continue
				}
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
				tag, tagErr := reflectCobraXFlag(field)
				if tagErr != nil {
					tag, _ = reflectLegacyFlag(field)
				}

				fieldTypeKind := fieldType.Kind()
				elementKind := reflect.Invalid
				if fieldTypeKind == reflect.Slice {
					elementKind = fieldType.Elem().Kind()
				}

				desc := NewFlagDescriptor(tag.Name, tag.Shorthand, tag.Usage, fieldTypeKind, elementKind, fieldValue)
				if tag.DefaultValue != "" && fieldValue.IsZero() {
					_ = desc.SetValueFromText(tag.DefaultValue)
				}
				flags = append(flags, desc)
			}
		}
	}

	return NewCommandDescriptor(use, shortHelpText, longHelpText, flags, arguments)
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
