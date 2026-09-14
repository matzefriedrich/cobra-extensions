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
func (r *commandReflector[T]) ReflectCommandDescriptor(handler T) types.CommandDescriptor {

	arguments := NewArgumentsDescriptorWith()

	value, valueType := resolveRootObject(handler)
	commandMetadata := extractCommandMetadataUse(valueType.Name())

	stack := utils.MakeStack[valueItem]()
	stack.Push(valueItem{value: value, valueType: valueType})

	flags := r.traverseFields(stack, commandMetadata, arguments)

	return NewCommandDescriptor(commandMetadata.use, commandMetadata.shortHelpText, commandMetadata.longHelpText, flags, arguments)
}

func (r *commandReflector[T]) traverseFields(stack utils.Stack[valueItem], commandMetadata *commandMetadata, arguments types.ArgumentsDescriptor) []FlagDescriptor {
	flags := make([]FlagDescriptor, 0)

	for !stack.IsEmpty() {

		next := stack.Pop()

		for i := range next.value.NumField() {

			field := next.valueType.Field(i)

			if shouldSkipCommandField(field.Type, field, commandMetadata) {
				continue
			}

			fieldValue := next.value.Field(i)

			if reflectArgumentsDescriptor(field.Type, fieldValue, arguments) {
				continue
			}

			if field.Anonymous {
				stack.Push(valueItem{value: fieldValue, valueType: field.Type})
				continue
			}

			if field.PkgPath == "" {
				flagDescriptor, ok := reflectFlagDescriptor(field, fieldValue)
				if ok {
					flags = append(flags, flagDescriptor)
				}
			}
		}
	}

	return flags
}

func resolveRootObject[T any](handler T) (reflect.Value, reflect.Type) {
	value := reflect.ValueOf(handler)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	valueType := reflect.TypeOf(value.Interface())
	return value, valueType
}

func extractCommandMetadataUse(valueTypeName string) *commandMetadata {
	return &commandMetadata{
		use: utils.ExtractCommandUse(valueTypeName),
	}
}

type commandMetadata struct {
	use           string
	shortHelpText string
	longHelpText  string
}

func (m *commandMetadata) Update(tag *CobraXCommandTag) {
	m.use = ternary.ValueOrDefault(tag.Use, ternary.NotNilOrWhitespace, m.use)
	m.shortHelpText = ternary.ValueOrDefault(tag.Help, ternary.NotNilOrWhitespace, m.shortHelpText)
	m.longHelpText = ternary.ValueOrDefault(tag.Description, ternary.NotNilOrWhitespace, m.longHelpText)
}

func isCommandFieldType(fieldType reflect.Type) bool {
	//nolint:staticcheck // required for legacy functionality
	return fieldType == reflect.TypeFor[types.CommandName]() || reflect.TypeFor[types.BaseCommand]() == fieldType
}

func isCommandNameFieldType(fieldType reflect.Type) bool {
	//nolint:staticcheck // required for legacy functionality
	return fieldType == reflect.TypeFor[types.CommandName]()
}

// shouldSkipCommandField resolves the command tag for command-type fields, applies it to the metadata,
// and reports whether the field should be skipped during reflection.
func shouldSkipCommandField(fieldType reflect.Type, field reflect.StructField, metadata *commandMetadata) bool {
	if !isCommandFieldType(fieldType) {
		return false
	}
	commandTag, skipField := resolveCommandTag(field)
	if skipField {
		return true
	}
	if commandTag != nil {
		metadata.Update(commandTag)
	}
	return isCommandNameFieldType(fieldType)
}

// resolveCommandTag resolves the command tag using the cobra-x convention and falls back to the legacy convention.
// skipField indicates that neither convention applies and the field should be ignored.
func resolveCommandTag(field reflect.StructField) (commandTag *CobraXCommandTag, skipField bool) {
	tag, tagErr := reflectCobraXCommand(field)
	if tagErr != nil && errors.Is(tagErr, ErrCobraXCommandNotFound) {
		tag, tagErr = reflectLegacyCommand(field)
		if tagErr != nil && errors.Is(tagErr, ErrCobraXLegacyTagsNotFound) {
			return nil, true
		}
	}
	return tag, false
}

func reflectFlagDescriptor(field reflect.StructField, fieldValue reflect.Value) (FlagDescriptor, bool) {
	tag, tagErr := reflectCobraXFlag(field)
	if tagErr != nil {
		tag, _ = reflectLegacyFlag(field)
	}

	if tag == nil {
		return FlagDescriptor{}, false
	}

	fieldTypeKind := field.Type.Kind()
	elementKind := reflect.Invalid
	if fieldTypeKind == reflect.Slice {
		elementKind = field.Type.Elem().Kind()
	}

	descriptor := NewFlagDescriptor(tag.Name, tag.Shorthand, tag.Usage, fieldTypeKind, elementKind, fieldValue)
	if tag.DefaultValue != "" && fieldValue.IsZero() {
		_ = descriptor.SetValueFromText(tag.DefaultValue)
	}
	if tag.SettingKey != "" {
		descriptor = descriptor.WithSettingKey(tag.SettingKey)
	}
	return descriptor, true
}

func reflectArgumentsDescriptor(fieldType reflect.Type, fieldValue reflect.Value, arguments types.ArgumentsDescriptor) bool {
	hasCommandArgs := false
	reflectedObject := ReflectedObject{instanceValue: fieldValue, objectType: fieldType}
	reflectedObject.EnumerateFields(func(index int, field ReflectedField) {
		fieldTypeKind := field.typeKind()
		switch fieldTypeKind {
		case reflect.String:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Bool:
			if hasCommandArgs {
				descriptor := ArgumentDescriptor{typeKind: fieldTypeKind, value: field.value, argumentIndex: index - 1, displayName: field.cobraXArgumentName()}
				arguments.With(Args(descriptor))
			}
		case reflect.Interface:
		case reflect.Struct:
			if field.isType(types.CommandArgs{}) {
				compatible, ok := field.getInterfaceValue().(types.CommandArgs)
				if ok {
					arguments.With(MinimumArgs(compatible.MinimumArgs))
					hasCommandArgs = true
				}
			}
		}
	})
	return hasCommandArgs
}
