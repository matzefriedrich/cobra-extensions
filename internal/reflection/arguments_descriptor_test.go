package reflection

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_argumentsDescriptor_bind_arguments_appends_positional_placeholders_to_use(t *testing.T) {
	tests := []struct {
		name        string
		minimumArgs int
		args        []ArgumentDescriptor
		baseUse     string
		expectedUse string
	}{
		{
			name:        "renders_required_placeholder_when_minimum_args_is_satisfied",
			minimumArgs: 1,
			args:        []ArgumentDescriptor{{argumentIndex: 0, displayName: "name"}},
			baseUse:     "greet",
			expectedUse: "greet <name>",
		},
		{
			name:        "renders_optional_placeholder_when_minimum_args_is_zero",
			minimumArgs: 0,
			args:        []ArgumentDescriptor{{argumentIndex: 0, displayName: "name"}},
			baseUse:     "greet",
			expectedUse: "greet [name]",
		},
		{
			name:        "renders_required_then_optional_in_declaration_order",
			minimumArgs: 1,
			args: []ArgumentDescriptor{
				{argumentIndex: 0, displayName: "name"},
				{argumentIndex: 1, displayName: "other"},
			},
			baseUse:     "greet",
			expectedUse: "greet <name> [other]",
		},
		{
			name:        "skips_arguments_without_a_display_name",
			minimumArgs: 0,
			args: []ArgumentDescriptor{
				{argumentIndex: 0, displayName: ""},
				{argumentIndex: 1, displayName: "name"},
			},
			baseUse:     "greet",
			expectedUse: "greet [name]",
		},
		{
			name:        "required_threshold_follows_argument_index",
			minimumArgs: 1,
			args: []ArgumentDescriptor{
				{argumentIndex: 0, displayName: ""},
				{argumentIndex: 1, displayName: "name"},
			},
			baseUse:     "greet",
			expectedUse: "greet [name]",
		},
		{
			name:        "leaves_use_unchanged_when_there_are_no_arguments",
			minimumArgs: 0,
			args:        nil,
			baseUse:     "greet",
			expectedUse: "greet",
		},
		{
			name:        "renders_placeholders_without_leading_space_when_use_is_empty",
			minimumArgs: 1,
			args:        []ArgumentDescriptor{{argumentIndex: 0, displayName: "name"}},
			baseUse:     "",
			expectedUse: "<name>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			arguments := NewArgumentsDescriptorWith(Args(tt.args...), MinimumArgs(tt.minimumArgs))
			target := &cobra.Command{Use: tt.baseUse}

			// Act
			arguments.BindArguments(target)

			// Assert
			assert.Equal(t, tt.expectedUse, target.Use)
		})
	}
}