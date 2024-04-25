package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"

	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

type helloCommand struct {
	use       abstractions.CommandName `flag:"hello" short:""`
	Arguments helloArgs
}

type helloArgs struct {
	abstractions.CommandArgs
	Name string
}

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
