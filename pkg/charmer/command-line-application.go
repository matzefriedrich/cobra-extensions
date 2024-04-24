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

func (a *CommandLineApplication) AddCommand(c *cobra.Command, setup CommandsSetupFunc) *CommandLineApplication {
	a.root.AddCommand(c)
	wrapper := commandSetup{command: c}
	if setup != nil {
		setup(&wrapper)
	}
	return a
}
