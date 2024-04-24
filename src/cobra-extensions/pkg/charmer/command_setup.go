package charmer

import "github.com/spf13/cobra"

type CommandsSetupFunc func(w CommandSetup)

type CommandSetup interface {
	AddCommand(c *cobra.Command, setup CommandsSetupFunc) CommandSetup
}

var _ CommandSetup = &commandSetup{}

type commandSetup struct {
	command *cobra.Command
}

// AddCommand Adds a sub-command to the specified *cobra.Command instance.
func (w *commandSetup) AddCommand(c *cobra.Command, setup CommandsSetupFunc) CommandSetup {
	w.command.AddCommand(c)
	wrapper := commandSetup{command: c}
	if setup != nil {
		setup(&wrapper)
	}
	return w
}
