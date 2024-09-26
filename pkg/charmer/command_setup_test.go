package charmer

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"slices"
	"testing"
)

func Test_CommandSetup_AddGroupCommand_adds_subcommand_via_setup_func(t *testing.T) {

	// Arrange
	root := &cobra.Command{}
	sut := newCommandSetup(root)

	// Act
	group := &cobra.Command{Use: "group"}
	c := &cobra.Command{Use: "command"}

	sut.AddGroupCommand(group, func(setup types.CommandSetup) {
		setup.AddCommand(c)
	})

	actual := slices.Contains(group.Commands(), c)

	// Assert
	assert.True(t, actual)
}
