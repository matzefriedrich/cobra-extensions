package reflection

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FlagDescriptor_set_value_returns_error_for_invalid_types_or_values(t *testing.T) {
	// Arrange
	sValue := ""
	iValue := int64(0)
	bValue := false
	slValue := []string{}

	tests := []struct {
		name      string
		kind      reflect.Kind
		elemKind  reflect.Kind
		target    reflect.Value
		value     any
		expectErr string
	}{
		{
			name:      "string field with int value",
			kind:      reflect.String,
			elemKind:  reflect.Invalid,
			target:    reflect.ValueOf(&sValue).Elem(),
			value:     123,
			expectErr: ErrorInvalidValue,
		},
		{
			name:      "int64 field with string value",
			kind:      reflect.Int64,
			elemKind:  reflect.Invalid,
			target:    reflect.ValueOf(&iValue).Elem(),
			value:     "not-int",
			expectErr: ErrorInvalidValue,
		},
		{
			name:      "bool field with string value",
			kind:      reflect.Bool,
			elemKind:  reflect.Invalid,
			target:    reflect.ValueOf(&bValue).Elem(),
			value:     "not-bool",
			expectErr: ErrorInvalidValue,
		},
		{
			name:      "slice field with non-slice value",
			kind:      reflect.Slice,
			elemKind:  reflect.String,
			target:    reflect.ValueOf(&slValue).Elem(),
			value:     "not-slice",
			expectErr: ErrorInvalidValue,
		},
		{
			name:      "unsupported float64 field",
			kind:      reflect.Float64,
			elemKind:  reflect.Invalid,
			target:    reflect.Value{},
			value:     1.23,
			expectErr: ErrorFlagTypeNotSupported,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := NewFlagDescriptor("test", "", "usage", tt.kind, tt.elemKind, tt.target)

			// Act
			err := d.SetValue(tt.value)

			// Assert
			assert.Error(t, err)
			assert.Equal(t, tt.expectErr, err.Error())
		})
	}
}
