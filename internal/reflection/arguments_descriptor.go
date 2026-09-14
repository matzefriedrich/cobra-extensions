package reflection

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
	"strings"
)

// ArgumentsDescriptor Stores arguments metadata.
type argumentsDescriptor struct {
	minimumArgs int
	args        []ArgumentDescriptor
}

var _ types.ArgumentsDescriptor = (*argumentsDescriptor)(nil)

// BindArguments sets the minimum number of positional arguments required for the given Cobra command.
func (d *argumentsDescriptor) BindArguments(target *cobra.Command) {
	target.Args = cobra.MinimumNArgs(d.minimumArgs)
	d.appendUseArgumentPlaceholders(target)
}

// appendUseArgumentPlaceholders appends the positional argument placeholders to the command's Use string.
func (d *argumentsDescriptor) appendUseArgumentPlaceholders(target *cobra.Command) {
	placeholders := d.renderPlaceholders()
	if len(placeholders) == 0 {
		return
	}
	useWithPlaceholders := target.Use + " " + strings.Join(placeholders, " ")
	target.Use = strings.TrimSpace(useWithPlaceholders)
}

func (d *argumentsDescriptor) renderPlaceholders() []string {
	placeholders := make([]string, 0, len(d.args))
	for _, argument := range d.args {
		placeholder := renderArgumentPlaceholder(argument, d.minimumArgs)
		if placeholder != "" {
			placeholders = append(placeholders, placeholder)
		}
	}
	return placeholders
}

func renderArgumentPlaceholder(argument ArgumentDescriptor, minimumArgs int) string {
	if argument.displayName == "" {
		return ""
	}
	if argument.argumentIndex < minimumArgs {
		return "<" + argument.displayName + ">"
	}
	return "[" + argument.displayName + "]"
}

// BindArgumentValues Sets the given set of values to positional argument fields.
func (d *argumentsDescriptor) BindArgumentValues(args ...string) {
	for _, argument := range d.args {
		if argument.argumentIndex >= len(args) {
			continue
		}
		value := args[argument.argumentIndex]
		argument.bindValue(value)
	}
}

func (a *ArgumentDescriptor) bindValue(value string) {
	switch a.typeKind {
	case reflect.String:
		a.SetString(value)
	case reflect.Int64:
		number, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			a.SetInt64(number)
		}
	case reflect.Bool:
		boolean, err := strconv.ParseBool(value)
		if err == nil {
			a.SetBool(boolean)
		}
	default:
		panic("unsupported type")
	}
}

// NewArgumentsDescriptorWith Creates a new ArgumentsDescriptor.
func NewArgumentsDescriptorWith(options ...types.ArgumentsDescriptorOption) types.ArgumentsDescriptor {
	argumentsDescriptor := &argumentsDescriptor{
		args: make([]ArgumentDescriptor, 0),
	}
	for _, option := range options {
		option(argumentsDescriptor)
	}
	return argumentsDescriptor
}
