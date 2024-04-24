package charmer

import "github.com/spf13/cobra"

type CommandLineApplication struct {
	root *cobra.Command
}

// NewCommandLineApplication Creates a new CommandLineApplication instance.
func NewCommandLineApplication() *CommandLineApplication {
	return &CommandLineApplication{
		root: &cobra.Command{},
	}
}

func (a *CommandLineApplication) Execute() error {
	return a.root.Execute()
}

func (a *CommandLineApplication) AddCommand(c ...*cobra.Command) *CommandLineApplication {
	a.root.AddCommand(c...)
	return a
}

func (a *CommandLineApplication) AddGroupCommand(c *cobra.Command, setup CommandsSetupFunc) *CommandLineApplication {
	a.root.AddCommand(c)
	if setup != nil {
		wrapper := commandSetup{command: c}
		setup(&wrapper)
	}
	return a
}
