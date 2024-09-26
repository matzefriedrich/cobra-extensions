package charmer

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CommandLineApplication_Execute_without_any_command_does_not_return_error(t *testing.T) {

	// Arrange
	sut := NewCommandLineApplication("test-app", "")

	// Act
	err := sut.Execute()

	// Assert
	assert.NoError(t, err)
}

func Test_CommandLineApplication_Execute_smoke_test(t *testing.T) {

	// Arrange
	application := NewCommandLineApplication("test-app", "")

	command := newTestCommand(noop())
	sut := command.AsTypedCommand()
	application.AddGroupCommand(newTestGroupCommand(), func(setup types.CommandSetup) {
		setup.AddCommand(sut)
	})

	application.root.SetArgs([]string{"group", "test"})

	// Act
	err := application.Execute()

	// Assert
	assert.NoError(t, err)
	assert.True(t, command.Executed())
}

type executeFunc func()

func noop() executeFunc {
	return func() {}
}

type testCommand struct {
	use          types.CommandName `flag:"test"`
	typedCommand types.TypedCommand
	executeFunc  executeFunc
	executed     bool
}

func (t *testCommand) Executed() bool {
	return t.executed
}

func (t *testCommand) Execute() {
	t.executed = true
	t.executeFunc()
}

var _ types.TypedCommand = (*testCommand)(nil)

func newTestCommand(execute executeFunc) *testCommand {
	return &testCommand{
		executeFunc: execute,
	}
}

func (t *testCommand) AsTypedCommand() *cobra.Command {
	typedCommand := commands.CreateTypedCommand(t)
	return typedCommand
}

type testGroupCommand struct {
	types.BaseCommand
	use types.CommandName `flag:"group"`
}

func newTestGroupCommand() *cobra.Command {
	command := &testGroupCommand{BaseCommand: types.BaseCommand{}}
	return commands.CreateTypedCommand(command)
}
