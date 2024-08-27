package charmer

import "github.com/spf13/cobra"

type CommandLineApplication struct {
	root *cobra.Command
}

// NewRootCommand Creates a new cobra.Command object to be used as the application root command.
func NewRootCommand(name string, description string) *cobra.Command {
	if len(name) == 0 {
		panic("name is required")
	}
	return &cobra.Command{
		Use:   name,
		Short: description,
	}
}

// NewCommandLineApplication Creates a new CommandLineApplication instance.
func NewCommandLineApplication(name string, description string) *CommandLineApplication {
	return &CommandLineApplication{
		root: NewRootCommand(name, description),
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
