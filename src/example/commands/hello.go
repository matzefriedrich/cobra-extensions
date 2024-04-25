package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"

	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

// helloCommand A command handler type for the hello command.
type helloCommand struct {
	use       abstractions.CommandName `flag:"hello"`
	Arguments helloArgs
}

// helloArgs Stores values for positional arguments of the hello command.
type helloArgs struct {
	abstractions.CommandArgs
	Name string
}

// CreateHelloCommand Creates a new helloCommand instance.
func CreateHelloCommand() *cobra.Command {
	instance := &helloCommand{
		Arguments: helloArgs{
			CommandArgs: abstractions.NewCommandArgs(1),
		}}
	return pkg.CreateTypedCommand(instance)
}

func (c *helloCommand) Execute() {
	fmt.Printf("Hello %s.\n", c.Arguments.Name)
}
