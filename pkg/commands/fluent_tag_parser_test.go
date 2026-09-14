package commands

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type fluentTagCommand struct {
	types.BaseCommand `cobra-x:"fluent-test.description('Fluent flow test command')"`
	Name              string `cobra-x:"--name.shorthand(n).usage(Name to greet).default(World)"`
}

func (c *fluentTagCommand) Execute(_ context.Context) {
}

func Test_CreateTypedCommandWithTagParser_reflects_fluent_tag_keys(t *testing.T) {
	// Arrange
	instance := &fluentTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"fluent-test", "--name", "Ada"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewFluentTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Ada", instance.Name)
}

func Test_CreateTypedCommandWithTagParser_applies_default_values_from_fluent_keys(t *testing.T) {
	// Arrange
	instance := &fluentTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"fluent-test"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewFluentTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "World", instance.Name)
}
