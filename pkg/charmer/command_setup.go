package charmer

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type commandSetup struct {
	command *cobra.Command
}

var _ types.CommandSetup = &commandSetup{}

// newCommandSetup initializes a new commandSetup instance with the provided Cobra command object.
func newCommandSetup(command *cobra.Command) *commandSetup {
	return &commandSetup{command: command}
}

// AddCommand adds one or more sub-commands to the current command.
func (w *commandSetup) AddCommand(c ...*cobra.Command) types.CommandSetup {
	w.command.AddCommand(c...)
	return w
}

// AddGroupCommand adds a sub-command to the current command and calls the setup function for additional configuration.
func (w *commandSetup) AddGroupCommand(c *cobra.Command, setup types.CommandsSetupFunc) types.CommandSetup {
	w.command.AddCommand(c)
	if setup != nil {
		wrapper := commandSetup{command: c}
		setup(&wrapper)
	}
	return w
}
