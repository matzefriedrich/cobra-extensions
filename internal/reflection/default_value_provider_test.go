package reflection

import (
	"reflect"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type mockDefaultValueProvider struct {
	values map[string]string
	err    error
}

func (m *mockDefaultValueProvider) GetValue(key string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if val, ok := m.values[key]; ok {
		return val, nil
	}
	return "", nil
}

func Test_commandDescriptor_UnmarshalFlagValues_uses_default_value_provider(t *testing.T) {
	// Arrange
	provider := &mockDefaultValueProvider{
		values: map[string]string{
			"my-setting": "provider-value",
		},
	}
	var targetValue string
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())
	flag = flag.WithSettingKey("my-setting")

	descriptor := &commandDescriptor{
		flags:                []FlagDescriptor{flag},
		defaultValueProvider: provider,
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "", "")

	// Act
	descriptor.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "provider-value", targetValue)
}

func Test_commandDescriptor_UnmarshalFlagValues_prefers_command_line_value(t *testing.T) {
	// Arrange
	provider := &mockDefaultValueProvider{
		values: map[string]string{
			"my-setting": "provider-value",
		},
	}
	var targetValue string
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())
	flag = flag.WithSettingKey("my-setting")

	descriptor := &commandDescriptor{
		flags:                []FlagDescriptor{flag},
		defaultValueProvider: provider,
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "", "")
	_ = cmd.Flags().Set("my-flag", "cli-value")

	// Act
	descriptor.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "cli-value", targetValue)
}

func Test_reflectCobraXFlag_parses_setting_key_tag(t *testing.T) {
	// Arrange
	type testStruct struct {
		MyField string `cobra-x:"my-flag, setting-key=my-setting"`
	}
	structType := reflect.TypeFor[testStruct]()
	field, _ := structType.FieldByName("MyField")

	// Act
	tag, err := reflectCobraXFlag(field)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "my-setting", tag.SettingKey)
}

func Test_commandReflector_ReflectCommandDescriptor_populates_setting_key(t *testing.T) {
	// Arrange
	type testCommand struct {
		types.CommandName `cobra-x:"test"`
		MyFlag            string `cobra-x:"my-flag, setting-key=my-setting"`
	}
	handler := &testCommand{}
	reflector := NewCommandReflector[*testCommand]()

	// Act
	descriptor := reflector.ReflectCommandDescriptor(handler)
	d := descriptor.(*commandDescriptor)

	// Assert
	assert.Len(t, d.flags, 1)
	assert.Equal(t, "my-setting", d.flags[0].settingKey)
}

func Test_commandDescriptor_UnmarshalFlagValues_falls_back_to_default_on_provider_nil(t *testing.T) {
	// Arrange
	provider := &mockDefaultValueProvider{
		values: map[string]string{}, // empty
	}
	var targetValue string
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())
	flag = flag.WithSettingKey("my-setting")

	descriptor := &commandDescriptor{
		flags:                []FlagDescriptor{flag},
		defaultValueProvider: provider,
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "tag-default", "")

	// Act
	descriptor.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "tag-default", targetValue)
}

func Test_commandDescriptor_UnmarshalFlagValues_falls_back_to_default_on_provider_error(t *testing.T) {
	// Arrange
	provider := &mockDefaultValueProvider{
		err: assert.AnError,
	}

	var targetValue string
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())
	flag = flag.WithSettingKey("my-setting")

	descriptor := &commandDescriptor{
		flags:                []FlagDescriptor{flag},
		defaultValueProvider: provider,
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "tag-default", "")

	// Act
	descriptor.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "tag-default", targetValue)
}

func Test_commandDescriptor_UnmarshalFlagValues_works_without_provider(t *testing.T) {
	// Arrange
	var targetValue string
	flag := NewFlagDescriptor("my-flag", "", "", reflect.String, reflect.Invalid, reflect.ValueOf(&targetValue).Elem())

	descriptor := &commandDescriptor{
		flags: []FlagDescriptor{flag},
	}

	cmd := &cobra.Command{}
	cmd.Flags().String("my-flag", "tag-default", "")

	// Act
	descriptor.UnmarshalFlagValues(cmd)

	// Assert
	assert.Equal(t, "tag-default", targetValue)
}
