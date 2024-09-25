package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/reflection"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"

	"github.com/spf13/cobra"
)

// helloCommand A command handler type for the hello command.
type helloCommand struct {
	use       reflection.CommandName `flag:"hello"`
	Arguments helloArgs
}

var _ types.TypedCommand = (*helloCommand)(nil)

// helloArgs Stores values for positional arguments of the hello command.
type helloArgs struct {
	types.CommandArgs
	Name string
}

// CreateHelloCommand Creates a new helloCommand instance.
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
