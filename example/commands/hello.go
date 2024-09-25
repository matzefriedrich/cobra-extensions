package commands

import (
	"fmt"

	"github.com/matzefriedrich/cobra-extensions/pkg/commands"

	"github.com/matzefriedrich/cobra-extensions/pkg/types"

	"github.com/spf13/cobra"
)

type helloCommand struct {
	use       types.CommandName `flag:"hello" short:"Prints a greeting to the specified name."`
	Arguments helloArgs
}

var _ types.TypedCommand = (*helloCommand)(nil)

type helloArgs struct {
	types.CommandArgs
	Name string
}

func CreateHelloCommand() *cobra.Command {
	instance := &helloCommand{
		Arguments: helloArgs{
			CommandArgs: types.NewCommandArgs(types.MinimumArgumentsRequired(1)),
		}}
	return commands.CreateTypedCommand(instance)
}

func (c *helloCommand) Execute() {
	fmt.Printf("Hello %s.\n", c.Arguments.Name)
}
