package charmer

import (
	"context"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_CommandLineApplication_NewRootCommand(t *testing.T) {

	// Arrange
	const invalidRootCommandName = ""

	actExpectPanic := func(act func()) (panicked bool) {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		act()
		return false
	}

	// Act
	act := func() { _ = NewRootCommand(invalidRootCommandName, "") }
	result := make(chan bool, 1)
	go func() {
		result <- actExpectPanic(act)
	}()

	// Assert
	assert.True(t, <-result, "expected panic in goroutine creating root command")
}

func Test_CommandLineApplication_Execute_without_context_without_any_command_does_not_return_error(t *testing.T) {

	// Arrange
	sut := NewCommandLineApplication("test-app", "")

	// Act
	err := sut.Execute(nil)

	// Assert
	assert.NoError(t, err)
}

func Test_CommandLineApplication_Execute_without_any_command_does_not_return_error(t *testing.T) {

	// Arrange
	sut := NewCommandLineApplication("test-app", "")

	// Act
	err := sut.Execute(t.Context())

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
	err := application.Execute(t.Context())

	// Assert
	assert.NoError(t, err)
	assert.True(t, command.Executed())
}

type executeFunc func(_ context.Context)

func noop() executeFunc {
	return func(context.Context) {}
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

func (t *testCommand) Execute(ctx context.Context) {
	t.executed = true
	t.executeFunc(ctx)
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

func Test_CommandLineApplication_Execute_with_compound_command(t *testing.T) {
	// Arrange
	application := NewCommandLineApplication("test-application", "")

	application.AddCommand(NewTestCompoundCommand())

	application.root.SetArgs([]string{"compound", "--g1.value", "a", "--g2.value", "b"})

	// Act
	err := application.Execute(t.Context()) // TODO: Implement correct handling of compound flag structs

	// Assert
	assert.NoError(t, err)
}

type testCompoundCommand struct {
	types.BaseCommand
	use    types.CommandName `flag:"compound"`
	Group1 flagGroup         `flag:"g1"`
	Group2 flagGroup         `flag:"g2"`
}

type flagGroup struct {
	Value string `flag:"value"`
}

var _ types.TypedCommand = (*testCompoundCommand)(nil)

func NewTestCompoundCommand() *cobra.Command {
	instance := &testCompoundCommand{
		BaseCommand: types.BaseCommand{},
	}
	return commands.CreateTypedCommand(instance)
}
