package commands

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type shorthandCommand struct {
	name   types.CommandName `flag:"short"`
	Filter string            `flag:"filter" shorthand:"f" usage:"Filter value"`
	Toggle bool              `flag:"toggle" shorthand:"t" usage:"Toggle value"`
	Tags   []string          `flag:"tag" shorthand:"T" usage:"Tag values"`
}

func (s *shorthandCommand) Execute(_ context.Context) {
}

func Test_CreateTypedCommand_with_shorthand_flag(t *testing.T) {
	// Arrange
	instance := &shorthandCommand{}
	app := &cobra.Command{}
	app.SetArgs([]string{"short", "-f", "my-filter", "-t", "-T", "tag1", "-T", "tag2"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	err := app.ExecuteContext(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "my-filter", instance.Filter)
	assert.True(t, instance.Toggle)
	assert.Equal(t, []string{"tag1", "tag2"}, instance.Tags)
}
