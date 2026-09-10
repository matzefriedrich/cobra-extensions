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
	var slValue []string

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

func Test_FlagDescriptor_set_value_from_text_sets_values_for_supported_types(t *testing.T) {
	// Arrange
	sValue := ""
	iValue := int(0)
	i64Value := int64(0)
	bValue := false
	var ssValue []string
	var siValue []int
	var si64Value []int64
	var sbValue []bool

	tests := []struct {
		name     string
		kind     reflect.Kind
		elemKind reflect.Kind
		target   reflect.Value
		text     string
		expected any
	}{
		{name: "string value", kind: reflect.String, elemKind: reflect.Invalid, target: reflect.ValueOf(&sValue).Elem(), text: "hello", expected: "hello"},
		{name: "int value", kind: reflect.Int, elemKind: reflect.Invalid, target: reflect.ValueOf(&iValue).Elem(), text: "123", expected: 123},
		{name: "int64 value", kind: reflect.Int64, elemKind: reflect.Invalid, target: reflect.ValueOf(&i64Value).Elem(), text: "456", expected: int64(456)},
		{name: "bool value", kind: reflect.Bool, elemKind: reflect.Invalid, target: reflect.ValueOf(&bValue).Elem(), text: "true", expected: true},
		{name: "string slice value", kind: reflect.Slice, elemKind: reflect.String, target: reflect.ValueOf(&ssValue).Elem(), text: "a, b,c", expected: []string{"a", " b", "c"}},
		{name: "int slice value", kind: reflect.Slice, elemKind: reflect.Int, target: reflect.ValueOf(&siValue).Elem(), text: " 1, 2,3 ", expected: []int{1, 2, 3}},
		{name: "int64 slice value", kind: reflect.Slice, elemKind: reflect.Int64, target: reflect.ValueOf(&si64Value).Elem(), text: "4,5, 6", expected: []int64{4, 5, 6}},
		{name: "bool slice value", kind: reflect.Slice, elemKind: reflect.Bool, target: reflect.ValueOf(&sbValue).Elem(), text: "true,false,true", expected: []bool{true, false, true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := NewFlagDescriptor("test", "", "usage", tt.kind, tt.elemKind, tt.target)

			// Act
			err := d.SetValueFromText(tt.text)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, tt.target.Interface())
		})
	}
}

func Test_FlagDescriptor_set_value_from_text_returns_error_for_invalid_text(t *testing.T) {
	// Arrange
	iValue := int(0)
	bValue := false
	siValue := []int{}
	si64Value := []int64{}
	sbValue := []bool{}

	tests := []struct {
		name     string
		kind     reflect.Kind
		elemKind reflect.Kind
		target   reflect.Value
		text     string
	}{
		{name: "invalid int", kind: reflect.Int, elemKind: reflect.Invalid, target: reflect.ValueOf(&iValue).Elem(), text: "not-an-int"},
		{name: "invalid bool", kind: reflect.Bool, elemKind: reflect.Invalid, target: reflect.ValueOf(&bValue).Elem(), text: "not-a-bool"},
		{name: "invalid int slice element", kind: reflect.Slice, elemKind: reflect.Int, target: reflect.ValueOf(&siValue).Elem(), text: "1,not-an-int"},
		{name: "invalid int64 slice element", kind: reflect.Slice, elemKind: reflect.Int64, target: reflect.ValueOf(&si64Value).Elem(), text: "3,not-an-int64"},
		{name: "invalid bool slice element", kind: reflect.Slice, elemKind: reflect.Bool, target: reflect.ValueOf(&sbValue).Elem(), text: "true,not-a-bool"},
		{name: "unsupported type", kind: reflect.Float64, elemKind: reflect.Invalid, target: reflect.Value{}, text: "1.23"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			d := NewFlagDescriptor("test", "", "usage", tt.kind, tt.elemKind, tt.target)

			// Act
			err := d.SetValueFromText(tt.text)

			// Assert
			assert.Error(t, err)
		})
	}
}
