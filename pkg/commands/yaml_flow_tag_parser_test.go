package commands

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type yamlFlowTagCommand struct {
	types.BaseCommand `cobra-x:"{name: 'yaml-test', description: 'Yaml flow test command'}"`
	Name              string `cobra-x:"{name: '--name', shorthand: 'n', usage: 'Name to greet', default: 'World'}"`
}

func (c *yamlFlowTagCommand) Execute(_ context.Context) {
}

func Test_CreateTypedCommandWithTagParser_reflects_yaml_flow_tag_keys(t *testing.T) {
	// Arrange
	instance := &yamlFlowTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"yaml-test", "--name", "Ada"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewYamlFlowTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Ada", instance.Name)
}

func Test_CreateTypedCommandWithTagParser_applies_default_values_from_yaml_flow_keys(t *testing.T) {
	// Arrange
	instance := &yamlFlowTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"yaml-test"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewYamlFlowTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "World", instance.Name)
}
