package pkg

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testCommand1 struct {
	name abstractions.CommandName `flag:"test"`
	P1   string                   `flag:"param1"`
}

func (t *testCommand1) Execute() {
}

func Test_CreateTypedCommand(t *testing.T) {

	// Arrange
	instance := &testCommand1{}
	app := &cobra.Command{}
	app.SetArgs([]string{"test", "--param1", "Hello World"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.Execute()

	// Assert
	assert.Equal(t, "Hello World", instance.P1)
}

type testCommand2 struct {
	testCommand1
	name        abstractions.CommandName `flag:"test"`
	P2          int64                    `flag:"param2"`
	P3          int                      `flag:"p3"`
	BooleanFlag bool                     `flag:"b"`
}

func Test_CreateTypedCommand_with_base_template(t *testing.T) {

	// Arrange
	instance := &testCommand2{testCommand1: testCommand1{}}
	app := &cobra.Command{}
	app.SetArgs([]string{"test", "--param1", "Hello World", "--param2", "265", "--p3", "345", "--b"})

	cmd := CreateTypedCommand(instance)
	app.AddCommand(cmd)

	// Act
	_ = app.Execute()

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
	_ = app.Execute()

	// Assert
	assert.Equal(t, expectedP1, instance.P1)
	assert.Equal(t, expectedP2, instance.P2)
	assert.Equal(t, 345, instance.P3)
	assert.False(t, instance.BooleanFlag)
}

type testCommandWithPositionalArgs struct {
	use       abstractions.CommandName `flag:"test3"`
	Arguments testCommandArgs
}

type testCommandArgs struct {
	abstractions.CommandArgs
	TextArgument    string
	NumericArgument int64
	BooleanArgument bool
}

func (t *testCommandWithPositionalArgs) Execute() {

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
	_ = app.Execute()

	// Assert
	arguments := instance.Arguments

	assert.Equal(t, "Hello", arguments.TextArgument)
	assert.Equal(t, int64(5), arguments.NumericArgument)
	assert.Equal(t, true, arguments.BooleanArgument)
}
