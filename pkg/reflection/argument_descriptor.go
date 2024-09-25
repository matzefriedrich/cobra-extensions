package reflection

import (
	"reflect"
)

// ArgumentDescriptor represents metadata for an argument, including its index, value, and type kind.
type ArgumentDescriptor struct {
	argumentIndex int
	value         reflect.Value
	typeKind      reflect.Kind
}

// ArgumentIndex returns the index of the argument represented by the ArgumentDescriptor.
func (a *ArgumentDescriptor) ArgumentIndex() int {
	return a.argumentIndex
}

// SetString Applies a string argument to the underlying structure.
func (a *ArgumentDescriptor) SetString(value string) {
	a.value.SetString(value)
}

// SetInt64 Applies an int64 argument to the underlying structure.
func (a *ArgumentDescriptor) SetInt64(n int64) {
	a.value.SetInt(n)
}

// SetBool Applies a boolean argument to the underlying structure.
func (a *ArgumentDescriptor) SetBool(b bool) {
	a.value.SetBool(b)
}
