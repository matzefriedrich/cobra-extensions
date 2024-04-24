package charmer

import "github.com/spf13/cobra"

type CommandsSetupFunc func(w CommandSetup)

type CommandSetup interface {
	AddCommand(c ...*cobra.Command) CommandSetup
	AddGroupCommand(c *cobra.Command, setup CommandsSetupFunc) CommandSetup
}

var _ CommandSetup = &commandSetup{}

type commandSetup struct {
	command *cobra.Command
}

func (w *commandSetup) AddCommand(c ...*cobra.Command) CommandSetup {
	w.command.AddCommand(c...)
	return w
}

// AddGroupCommand Adds a sub-command to the specified *cobra.Command instance.
func (w *commandSetup) AddGroupCommand(c *cobra.Command, setup CommandsSetupFunc) CommandSetup {
	w.command.AddCommand(c)
	if setup != nil {
		wrapper := commandSetup{command: c}
		setup(&wrapper)
	}
	return w
}
