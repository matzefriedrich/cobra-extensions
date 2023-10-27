package commands

import (
	"fmt"

	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

type helloCommand struct {
	use  pkg.CommandName `flag:"hello"`
	Name string          `flag:"name" usage:"Your name"`
}

func CreateHelloCommand() *cobra.Command {
	instance := &helloCommand{}
	return pkg.CreateTypedCommand(instance)
}

func (c *helloCommand) Execute() {
	_ = fmt.Sprintf("Hello %s.", c.Name)
}
