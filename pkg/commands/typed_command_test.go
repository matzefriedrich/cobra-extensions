package commands

import (
	"context"
	"io"
	"testing"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testCommand1 struct {
	name types.CommandName `flag:"test"`
	P1   string            `flag:"param1"`
}

func (t *testCommand1) Execute(_ context.Context) {
}

func Test_CreateTypedCommand(t *testing.T) {

	// Arrange
	instance := &testCommand1{}
	app := &cobra.Command{}
	app.SetArgs([]string{"test", "--param1", "Hello World"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.ExecuteContext(t.Context())

	// Assert
	assert.Equal(t, "Hello World", instance.P1)
}

type logger struct {
	w                    io.Writer
	unsupportedFieldType int8
}

type testCommand2 struct {
	testCommand1
	name        types.CommandName `flag:"test"`
	P2          int64             `flag:"param2"`
	P3          int               `flag:"p3"`
	BooleanFlag bool              `flag:"b"`
	logger      logger
}

func Test_CreateTypedCommand_with_base_template(t *testing.T) {

	// Arrange
	instance := &testCommand2{testCommand1: testCommand1{}}
	app := &cobra.Command{}
	app.SetArgs([]string{"test", "--param1", "Hello World", "--param2", "265", "--p3", "345", "--b"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.ExecuteContext(t.Context())

	// Assert
	assert.Equal(t, "Hello World", instance.P1)
	assert.Equal(t, int64(265), instance.P2)
	assert.Equal(t, 345, instance.P3)
	assert.True(t, instance.BooleanFlag)
}

func Test_CreateTypedCommand_with_base_template_default_values(t *testing.T) {

	// Arrange
	expectedP1 := "not set"
	expectedP2 := int64(76)

	instance := &testCommand2{
		testCommand1: testCommand1{
			P1: expectedP1,
		},
		P2: expectedP2,
	}

	app := &cobra.Command{}
	app.SetArgs([]string{"test", "--p3", "345"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.ExecuteContext(t.Context())

	// Assert
	assert.Equal(t, expectedP1, instance.P1)
	assert.Equal(t, expectedP2, instance.P2)
	assert.Equal(t, 345, instance.P3)
	assert.False(t, instance.BooleanFlag)
}

type testCommandWithPositionalArgs struct {
	use       types.CommandName `flag:"test3"`
	Arguments testCommandArgs
}

type testCommandArgs struct {
	types.CommandArgs
	TextArgument    string
	NumericArgument int64
	BooleanArgument bool
}

func (t *testCommandWithPositionalArgs) Execute(_ context.Context) {

}

func Test_CreateTypedCommand_with_positional_args(t *testing.T) {

	// Arrange
	instance := &testCommandWithPositionalArgs{
		Arguments: testCommandArgs{},
	}

	app := &cobra.Command{}
	app.SetArgs([]string{"test3", "Hello", "5", "true"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.ExecuteContext(t.Context())

	// Assert
	arguments := instance.Arguments

	assert.Equal(t, "Hello", arguments.TextArgument)
	assert.Equal(t, int64(5), arguments.NumericArgument)
	assert.Equal(t, true, arguments.BooleanArgument)
}

func Test_NonRunnable_unsets_Run_and_RunE_fields(t *testing.T) {
	// Arrange
	sut := NonRunnable()
	target := &cobra.Command{}

	// Act
	sut(target)

	// Assert
	assert.Nilf(t, target.Run, "Command should not have a Run function")
	assert.Nilf(t, target.RunE, "Command should not have a RunE function")
}

type disabledCommand struct {
	execute func()
}

func (d *disabledCommand) Execute(_ context.Context) {
	d.execute()
}

var _ types.TypedCommand = (*disabledCommand)(nil)

func Test_CreateTypedCommand_with_NonRunnable_disables_the_Execute_handler(t *testing.T) {
	// Arrange
	executed := false
	instance := &disabledCommand{
		execute: func() {
			executed = true
		},
	}
	sut := CreateTypedCommand(instance, NonRunnable)

	// Act
	sut.ExecuteContext(t.Context())

	// Assert
	assert.NotNil(t, sut)
	assert.False(t, executed)
}
