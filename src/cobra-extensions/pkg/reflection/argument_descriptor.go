package reflection

import (
	"github.com/spf13/cobra"
	"reflect"
)

// ArgumentDescriptor Stores the position and a value reference for single positional argument.
type ArgumentDescriptor struct {
	argumentIndex int
	value         reflect.Value
}

// ArgumentIndex Gets a value indicating the argument position in the args array.
func (a *ArgumentDescriptor) ArgumentIndex() int {
	return a.argumentIndex
}

// SetValue Applies the argument value to the underlying structure.
func (a *ArgumentDescriptor) SetValue(value string) error {
	a.value.SetString(value)
	return nil
}

type ArgumentsDescriptor interface {
	With(options ...ArgumentsDescriptorOption) ArgumentsDescriptor
	BindArguments(target *cobra.Command)
	BindArgumentValues(args ...string)
}

// ArgumentsDescriptor Stores arguments metadata.
type argumentsDescriptor struct {
	minimumArgs int
	args        []ArgumentDescriptor
}

func (d *argumentsDescriptor) BindArguments(target *cobra.Command) {
	target.Args = cobra.MinimumNArgs(d.minimumArgs)
}

// BindArgumentValues Sets the given set of values to positional argument fields.
func (d *argumentsDescriptor) BindArgumentValues(args ...string) {
	for _, a := range d.args {
		index := a.argumentIndex
		if index < len(args) {
			value := args[index]
			_ = a.SetValue(value)
		}
	}
}

// NewArgumentsDescriptorWith Creates a new ArgumentsDescriptor.
func NewArgumentsDescriptorWith(options ...ArgumentsDescriptorOption) ArgumentsDescriptor {
	argumentsDescriptor := &argumentsDescriptor{
		args: make([]ArgumentDescriptor, 0),
	}
	for _, option := range options {
		option.Apply(argumentsDescriptor)
	}
	return argumentsDescriptor
}

type ArgumentsDescriptorOption interface {
	Apply(argumentsDescriptor *argumentsDescriptor)
}

type argumentsDescriptorOption struct {
	f func(argumentsDescriptor *argumentsDescriptor)
}

func newOption(f func(argumentsDescriptor *argumentsDescriptor)) ArgumentsDescriptorOption {
	return &argumentsDescriptorOption{f: f}
}

func (o *argumentsDescriptorOption) Apply(argumentsDescriptor *argumentsDescriptor) {
	o.f(argumentsDescriptor)
}

func MinimumArgs(value int) ArgumentsDescriptorOption {
	return newOption(func(descriptor *argumentsDescriptor) {
		descriptor.minimumArgs = value
	})
}

func Args(args ...ArgumentDescriptor) ArgumentsDescriptorOption {
	return newOption(func(descriptor *argumentsDescriptor) {
		descriptor.args = append(descriptor.args, args...)
	})
}

func (d *argumentsDescriptor) With(options ...ArgumentsDescriptorOption) ArgumentsDescriptor {
	for _, option := range options {
		option.Apply(d)
	}
	return d
}
