package reflection

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_commandDescriptor_bind_flags_correctly_maps_flags_to_cobra_command(t *testing.T) {
	// Arrange
	sValue := ""
	iValue := int(0)
	i64Value := int64(0)
	bValue := false
	ssValue := []string{}
	siValue := []int{}
	si64Value := []int64{}
	sbValue := []bool{}

	tests := []struct {
		name      string
		flagName  string
		shorthand string
		kind      reflect.Kind
		elemKind  reflect.Kind
		target    reflect.Value
	}{
		{name: "string with shorthand", flagName: "string-s", shorthand: "s", kind: reflect.String, elemKind: reflect.Invalid, target: reflect.ValueOf(&sValue).Elem()},
		{name: "string without shorthand", flagName: "string", shorthand: "", kind: reflect.String, elemKind: reflect.Invalid, target: reflect.ValueOf(&sValue).Elem()},
		{name: "int with shorthand", flagName: "int-i", shorthand: "i", kind: reflect.Int, elemKind: reflect.Invalid, target: reflect.ValueOf(&iValue).Elem()},
		{name: "int without shorthand", flagName: "int", shorthand: "", kind: reflect.Int, elemKind: reflect.Invalid, target: reflect.ValueOf(&iValue).Elem()},
		{name: "int64 with shorthand", flagName: "int64-i", shorthand: "j", kind: reflect.Int64, elemKind: reflect.Invalid, target: reflect.ValueOf(&i64Value).Elem()},
		{name: "int64 without shorthand", flagName: "int64", shorthand: "", kind: reflect.Int64, elemKind: reflect.Invalid, target: reflect.ValueOf(&i64Value).Elem()},
		{name: "bool with shorthand", flagName: "bool-b", shorthand: "b", kind: reflect.Bool, elemKind: reflect.Invalid, target: reflect.ValueOf(&bValue).Elem()},
		{name: "bool without shorthand", flagName: "bool", shorthand: "", kind: reflect.Bool, elemKind: reflect.Invalid, target: reflect.ValueOf(&bValue).Elem()},
		{name: "string slice with shorthand", flagName: "string-slice-S", shorthand: "S", kind: reflect.Slice, elemKind: reflect.String, target: reflect.ValueOf(&ssValue).Elem()},
		{name: "string slice without shorthand", flagName: "string-slice", shorthand: "", kind: reflect.Slice, elemKind: reflect.String, target: reflect.ValueOf(&ssValue).Elem()},
		{name: "int slice with shorthand", flagName: "int-slice-I", shorthand: "k", kind: reflect.Slice, elemKind: reflect.Int, target: reflect.ValueOf(&siValue).Elem()},
		{name: "int slice without shorthand", flagName: "int-slice", shorthand: "", kind: reflect.Slice, elemKind: reflect.Int, target: reflect.ValueOf(&siValue).Elem()},
		{name: "int64 slice with shorthand", flagName: "int64-slice-I", shorthand: "I", kind: reflect.Slice, elemKind: reflect.Int64, target: reflect.ValueOf(&si64Value).Elem()},
		{name: "int64 slice without shorthand", flagName: "int64-slice", shorthand: "", kind: reflect.Slice, elemKind: reflect.Int64, target: reflect.ValueOf(&si64Value).Elem()},
		{name: "bool slice with shorthand", flagName: "bool-slice-B", shorthand: "B", kind: reflect.Slice, elemKind: reflect.Bool, target: reflect.ValueOf(&sbValue).Elem()},
		{name: "bool slice without shorthand", flagName: "bool-slice", shorthand: "", kind: reflect.Slice, elemKind: reflect.Bool, target: reflect.ValueOf(&sbValue).Elem()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			flag := NewFlagDescriptor(tt.flagName, tt.shorthand, "usage", tt.kind, tt.elemKind, tt.target)
			descriptor := NewCommandDescriptor("test", "short", "long", []FlagDescriptor{flag}, NewArgumentsDescriptorWith())
			cmd := &cobra.Command{}

			// Act
			descriptor.BindFlags(cmd)

			// Assert
			f := cmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, f)
			assert.Equal(t, tt.shorthand, f.Shorthand)
		})
	}
}

func Test_commandDescriptor_unmarshal_flag_values_sets_field_values_from_cobra_command(t *testing.T) {
	// Arrange
	sValue := ""
	iValue := int(0)
	i64Value := int64(0)
	bValue := false
	ssValue := []string{}
	siValue := []int{}
	si64Value := []int64{}
	sbValue := []bool{}

	tests := []struct {
		name     string
		flagName string
		args     []string
		kind     reflect.Kind
		elemKind reflect.Kind
		target   reflect.Value
		expected any
	}{
		{name: "string flag", flagName: "string", args: []string{"--string", "hello"}, kind: reflect.String, elemKind: reflect.Invalid, target: reflect.ValueOf(&sValue).Elem(), expected: "hello"},
		{name: "int flag", flagName: "int", args: []string{"--int", "123"}, kind: reflect.Int, elemKind: reflect.Invalid, target: reflect.ValueOf(&iValue).Elem(), expected: 123},
		{name: "int64 flag", flagName: "int64", args: []string{"--int64", "456"}, kind: reflect.Int64, elemKind: reflect.Invalid, target: reflect.ValueOf(&i64Value).Elem(), expected: int64(456)},
		{name: "bool flag", flagName: "bool", args: []string{"--bool"}, kind: reflect.Bool, elemKind: reflect.Invalid, target: reflect.ValueOf(&bValue).Elem(), expected: true},
		{name: "string slice flag", flagName: "string-slice", args: []string{"--string-slice", "a", "--string-slice", "b"}, kind: reflect.Slice, elemKind: reflect.String, target: reflect.ValueOf(&ssValue).Elem(), expected: []string{"a", "b"}},
		{name: "int slice flag", flagName: "int-slice", args: []string{"--int-slice", "1", "--int-slice", "2"}, kind: reflect.Slice, elemKind: reflect.Int, target: reflect.ValueOf(&siValue).Elem(), expected: []int{1, 2}},
		{name: "int64 slice flag", flagName: "int64-slice", args: []string{"--int64-slice", "3", "--int64-slice", "4"}, kind: reflect.Slice, elemKind: reflect.Int64, target: reflect.ValueOf(&si64Value).Elem(), expected: []int64{3, 4}},
		{name: "bool slice flag", flagName: "bool-slice", args: []string{"--bool-slice", "true", "--bool-slice", "false"}, kind: reflect.Slice, elemKind: reflect.Bool, target: reflect.ValueOf(&sbValue).Elem(), expected: []bool{true, false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			flag := NewFlagDescriptor(tt.flagName, "", "usage", tt.kind, tt.elemKind, tt.target)
			descriptor := NewCommandDescriptor("test", "short", "long", []FlagDescriptor{flag}, NewArgumentsDescriptorWith())
			cmd := &cobra.Command{}
			descriptor.BindFlags(cmd)
			err := cmd.ParseFlags(tt.args)
			assert.NoError(t, err)

			// Act
			descriptor.UnmarshalFlagValues(cmd)

			// Assert
			assert.Equal(t, tt.expected, tt.target.Interface())
		})
	}
}

func Test_commandDescriptor_bind_arguments_sets_usage_descriptions(t *testing.T) {
	// Arrange
	argsDesc := NewArgumentsDescriptorWith(MinimumArgs(3))
	descriptor := NewCommandDescriptor("test-use", "test-short", "test-long", nil, argsDesc)
	cmd := &cobra.Command{}

	// Act
	descriptor.BindArguments(cmd)

	// Assert
	assert.Equal(t, "test-use", cmd.Use)
	assert.Equal(t, "test-short", cmd.Short)
	assert.Equal(t, "test-long", cmd.Long)
	// Verify minimum args (it's hard to check directly on cobra.Command because it's a function,
	// but we can try to call it and see if it fails)
	err := cmd.Args(cmd, []string{"a", "b"})
	assert.Error(t, err)
	err = cmd.Args(cmd, []string{"a", "b", "c"})
	assert.NoError(t, err)
}

func Test_argumentsDescriptor_with_applies_options(t *testing.T) {
	// Arrange
	argsDesc := NewArgumentsDescriptorWith()

	// Act
	argsDesc.With(MinimumArgs(5))

	// Assert
	cmd := &cobra.Command{}
	argsDesc.BindArguments(cmd)
	err := cmd.Args(cmd, []string{"1", "2", "3", "4"})
	assert.Error(t, err)
	err = cmd.Args(cmd, []string{"1", "2", "3", "4", "5"})
	assert.NoError(t, err)
}

func Test_commandDescriptor_unmarshal_argument_values_sets_positional_arguments(t *testing.T) {
	// Arrange
	sValue := ""
	i64Value := int64(0)
	bValue := false

	tests := []struct {
		name     string
		kind     reflect.Kind
		target   reflect.Value
		input    string
		expected any
	}{
		{name: "string argument", kind: reflect.String, target: reflect.ValueOf(&sValue).Elem(), input: "hello", expected: "hello"},
		{name: "int64 argument", kind: reflect.Int64, target: reflect.ValueOf(&i64Value).Elem(), input: "123", expected: int64(123)},
		{name: "bool argument", kind: reflect.Bool, target: reflect.ValueOf(&bValue).Elem(), input: "true", expected: true},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			argDesc := ArgumentDescriptor{argumentIndex: 0, value: tt.target, typeKind: tt.kind}
			argsDesc := NewArgumentsDescriptorWith(Args(argDesc))
			descriptor := NewCommandDescriptor("test", "short", "long", nil, argsDesc)

			// Act
			descriptor.UnmarshalArgumentValues(tt.input)

			// Assert
			assert.Equal(t, tt.expected, tt.target.Interface(), "Case %d: %s", i, tt.name)
		})
	}

	t.Run("invalid values do not change original values", func(t *testing.T) {
		// Arrange
		i64Value = 123
		bValue = true
		arg1 := ArgumentDescriptor{argumentIndex: 0, value: reflect.ValueOf(&i64Value).Elem(), typeKind: reflect.Int64}
		arg2 := ArgumentDescriptor{argumentIndex: 1, value: reflect.ValueOf(&bValue).Elem(), typeKind: reflect.Bool}
		argsDesc := NewArgumentsDescriptorWith(Args(arg1, arg2))
		descriptor := NewCommandDescriptor("test", "short", "long", nil, argsDesc)

		// Act
		descriptor.UnmarshalArgumentValues("not-an-int", "not-a-bool")

		// Assert
		assert.Equal(t, int64(123), i64Value)
		assert.True(t, bValue)
	})
}

func Test_argumentsDescriptor_bind_argument_values_panics_on_unsupported_type(t *testing.T) {
	// Arrange
	arg := 1.23
	desc := ArgumentDescriptor{argumentIndex: 0, value: reflect.ValueOf(&arg).Elem(), typeKind: reflect.Float64}
	argsDesc := NewArgumentsDescriptorWith(Args(desc))

	// Act & Assert
	assert.Panics(t, func() {
		argsDesc.BindArgumentValues("1.23")
	})
}

func Test_commandDescriptor_bind_flags_does_nothing_when_target_is_nil(t *testing.T) {
	descriptor := NewCommandDescriptor("test", "short", "long", nil, NewArgumentsDescriptorWith())
	assert.NotPanics(t, func() {
		descriptor.BindFlags(nil)
	})
}

func Test_commandDescriptor_bind_arguments_does_nothing_when_target_is_nil(t *testing.T) {
	descriptor := NewCommandDescriptor("test", "short", "long", nil, NewArgumentsDescriptorWith())
	assert.NotPanics(t, func() {
		descriptor.BindArguments(nil)
	})
}
