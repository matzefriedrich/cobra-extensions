package reflection

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"reflect"
	"strconv"
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
}

// BindArgumentValues Sets the given set of values to positional argument fields.
func (d *argumentsDescriptor) BindArgumentValues(args ...string) {
	for _, a := range d.args {
		index := a.argumentIndex
		if index < len(args) {
			value := args[index]
			switch a.typeKind {
			case reflect.String:
				a.SetString(value)
			case reflect.Int64:
				n, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					a.SetInt64(n)
				}
			case reflect.Bool:
				b, err := strconv.ParseBool(value)
				if err == nil {
					a.SetBool(b)
				}
			default:
				panic("unsupported type")
			}
		}
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
