package reflection

import "github.com/matzefriedrich/cobra-extensions/pkg/types"

func MinimumArgs(value int) types.ArgumentsDescriptorOption {
	return func(target any) {
		descriptor, ok := target.(*argumentsDescriptor)
		if ok {
			descriptor.minimumArgs = value
		}
	}
}

func Args(args ...ArgumentDescriptor) types.ArgumentsDescriptorOption {
	return func(target any) {
		descriptor, ok := target.(*argumentsDescriptor)
		if ok {
			descriptor.args = append(descriptor.args, args...)
		}
	}
}

func (d *argumentsDescriptor) With(options ...types.ArgumentsDescriptorOption) types.ArgumentsDescriptor {
	for _, option := range options {
		option(d)
	}
	return d
}
