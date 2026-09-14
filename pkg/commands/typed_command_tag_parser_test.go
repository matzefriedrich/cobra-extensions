package commands

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type standardTagCommand struct {
	types.BaseCommand `cobra-x:"standard-test" cobra-x-description:"Standard tag test command"`
	Name              string `cobra-x:"--name" cobra-x-shorthand:"n" cobra-x-usage:"Name to greet" cobra-x-default:"World"`
}

func (c *standardTagCommand) Execute(_ context.Context) {
}

func Test_CreateTypedCommandWithTagParser_reflects_standard_tag_keys(t *testing.T) {
	// Arrange
	instance := &standardTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"standard-test", "--name", "Ada"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewStandardTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "Ada", instance.Name)
}

func Test_CreateTypedCommandWithTagParser_applies_default_values_from_standard_keys(t *testing.T) {
	// Arrange
	instance := &standardTagCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"standard-test"})

	sut := CreateTypedCommandWithTagParser(instance, types.NewStandardTagParser())
	app.AddCommand(sut)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "World", instance.Name)
}
