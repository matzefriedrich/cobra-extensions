package reflection

import (
	"reflect"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

type behaviorCommand struct {
	types.BaseCommand `cobra-x:"behavior, help='short help', description='long help'"`
	UntaggedField     string
	StringFlag        string   `cobra-x:"--str, default='dflt', help='str usage'"`
	StringsFlag       []string `cobra-x:"--strs"`
	IntsFlag          []int    `cobra-x:"--ints, default='1,2'"`
	Int64sFlag        []int64  `cobra-x:"--i64s"`
	BoolsFlag         []bool   `cobra-x:"--bools"`
	BoolFlag          bool     `cobra-x:"--on"`
}

func Test_commandReflector_reflect_command_descriptor_creates_cobra_flags_with_expected_types_and_defaults(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*behaviorCommand]()
	handler := &behaviorCommand{}

	tests := []struct {
		name         string
		flagName     string
		expectedType string
		defaultValue string
	}{
		{name: "string flag", flagName: "str", expectedType: "string", defaultValue: "dflt"},
		{name: "string slice flag", flagName: "strs", expectedType: "stringSlice", defaultValue: "[]"},
		{name: "int slice flag", flagName: "ints", expectedType: "intSlice", defaultValue: "[1,2]"},
		{name: "int64 slice flag", flagName: "i64s", expectedType: "int64Slice", defaultValue: "[]"},
		{name: "bool slice flag", flagName: "bools", expectedType: "boolSlice", defaultValue: "[]"},
		{name: "bool flag", flagName: "on", expectedType: "bool", defaultValue: "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			descriptor := reflector.ReflectCommandDescriptor(handler)
			cmd := &cobra.Command{}
			descriptor.BindFlags(cmd)

			// Assert
			flag := cmd.Flags().Lookup(tt.flagName)
			assert.NotNil(t, flag)
			assert.Equal(t, tt.expectedType, flag.Value.Type())
			assert.Equal(t, tt.defaultValue, flag.DefValue)
		})
	}
}

func Test_commandReflector_reflect_command_descriptor_skips_untagged_fields(t *testing.T) {
	// Arrange
	type untaggedBaseCommand struct {
		types.BaseCommand
		Name string `cobra-x:"--name"`
	}
	type untaggedFieldsCommand struct {
		types.BaseCommand `cobra-x:"untagged"`
		Ignored           string
		Flag              string `cobra-x:"--flag"`
	}

	tests := []struct {
		name               string
		handler            any
		expectedUse        string
		expectedFlagNames  []string
		unexpectedFlagName string
	}{
		{
			name:               "embedded BaseCommand without tags is skipped and use is derived from type name",
			handler:            &untaggedBaseCommand{},
			expectedUse:        "untaggedbase",
			expectedFlagNames:  []string{"name"},
			unexpectedFlagName: "BaseCommand",
		},
		{
			name:               "exported field without tags yields no flag",
			handler:            &untaggedFieldsCommand{},
			expectedUse:        "untagged",
			expectedFlagNames:  []string{"flag"},
			unexpectedFlagName: "Ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			reflector := NewCommandReflector[any]()
			descriptor := reflector.ReflectCommandDescriptor(tt.handler)
			cmd := &cobra.Command{}
			descriptor.BindArguments(cmd)
			descriptor.BindFlags(cmd)

			// Assert
			assert.Equal(t, tt.expectedUse, cmd.Use)
			assert.Nil(t, cmd.Flags().Lookup(tt.unexpectedFlagName))
			assert.Equal(t, len(tt.expectedFlagNames), countDefinedFlags(cmd))
			for _, flagName := range tt.expectedFlagNames {
				assert.NotNil(t, cmd.Flags().Lookup(flagName))
			}
			assertContainsExactFlags(t, cmd, tt.expectedFlagNames)
		})
	}
}

func countDefinedFlags(cmd *cobra.Command) int {
	count := 0
	cmd.Flags().VisitAll(func(_ *pflag.Flag) {
		count++
	})
	return count
}

func assertContainsExactFlags(t *testing.T, cmd *cobra.Command, expected []string) {
	t.Helper()
	actual := make([]string, 0)
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		actual = append(actual, flag.Name)
	})
	assert.ElementsMatch(t, expected, actual)
}

func Test_commandReflector_reflect_command_descriptor_and_bind_pipeline_preserves_values(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*behaviorCommand]()

	tests := []struct {
		name            string
		args            []string
		expectedString  string
		expectedStrings []string
		expectedInts    []int
		expectedBool    bool
	}{
		{
			name:            "all flags parsed from command line",
			args:            []string{"--str", "cli", "--strs", "a", "--strs", "b", "--ints", "3", "--ints", "4", "--on"},
			expectedString:  "cli",
			expectedStrings: []string{"a", "b"},
			expectedInts:    []int{3, 4},
			expectedBool:    true,
		},
		{
			name:            "defaults applied when no flags are set",
			args:            []string{},
			expectedString:  "dflt",
			expectedStrings: []string{},
			expectedInts:    []int{1, 2},
			expectedBool:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			handler := &behaviorCommand{}

			// Act
			descriptor := reflector.ReflectCommandDescriptor(handler)
			cmd := &cobra.Command{}
			descriptor.BindFlags(cmd)
			err := cmd.ParseFlags(tt.args)
			descriptor.UnmarshalFlagValues(cmd)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedString, handler.StringFlag)
			assert.Equal(t, tt.expectedStrings, handler.StringsFlag)
			assert.Equal(t, tt.expectedInts, handler.IntsFlag)
			assert.Equal(t, tt.expectedBool, handler.BoolFlag)
		})
	}
}

func Test_commandReflector_reflect_command_descriptor_reflects_slice_element_kinds(t *testing.T) {
	// Arrange
	reflector := NewCommandReflector[*behaviorCommand]()
	handler := &behaviorCommand{}
	expectedElementKinds := map[string]reflect.Kind{
		"str":   reflect.Invalid,
		"strs":  reflect.String,
		"ints":  reflect.Int,
		"i64s":  reflect.Int64,
		"bools": reflect.Bool,
		"on":    reflect.Invalid,
	}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(handler).(*commandDescriptor)

	// Assert
	for _, flag := range descriptor.flags {
		assert.Equal(t, expectedElementKinds[flag.name], flag.elementKind, "flag %s", flag.name)
	}
	assert.Equal(t, len(expectedElementKinds), len(descriptor.flags))
}

func Test_commandReflector_reflect_command_descriptor_applies_slice_default_value_to_zero_value_fields(t *testing.T) {
	// Arrange
	type sliceDefaultsCommand struct {
		types.BaseCommand `cobra-x:"defaults"`
		StringSlice       []string `cobra-x:"--strs, default='x,y'"`
		Int64Slice        []int64  `cobra-x:"--i64s, default='7,8'"`
	}
	reflector := NewCommandReflector[*sliceDefaultsCommand]()
	handler := &sliceDefaultsCommand{}

	// Act
	descriptor := reflector.ReflectCommandDescriptor(handler)
	cmd := &cobra.Command{}
	descriptor.BindFlags(cmd)

	// Assert
	assert.Equal(t, "[x,y]", cmd.Flags().Lookup("strs").Value.String())
	assert.Equal(t, "[7,8]", cmd.Flags().Lookup("i64s").Value.String())
}

func Test_commandDescriptor_set_default_value_provider_then_unmarshal_uses_provider(t *testing.T) {
	// Arrange
	provider := &mockDefaultValueProvider{
		values: map[string]string{"my-setting": "provider-value"},
	}
	var targetValue string
	desc := NewCommandDescriptor("test", "short", "long", []FlagDescriptor{}, NewArgumentsDescriptorWith())
	defaultDesc := desc.(*commandDescriptor)
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())
	flag = flag.WithSettingKey("my-setting")
	defaultDesc.flags = append(defaultDesc.flags, flag)

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "", "")

	// Act
	defaultDesc.SetDefaultValueProvider(provider)
	defaultDesc.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "provider-value", targetValue)
}