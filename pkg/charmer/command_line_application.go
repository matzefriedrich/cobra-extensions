package charmer

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

// CommandLineApplication Represents a command-line application using the Cobra library for command parsing.
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

// Execute executes the root command of the CommandLineApplication.
func (a *CommandLineApplication) Execute() error {
	return a.root.Execute()
}

// AddCommand Adds one or more commands to the root command of the CommandLineApplication.
func (a *CommandLineApplication) AddCommand(c ...*cobra.Command) *CommandLineApplication {
	a.root.AddCommand(c...)
	return a
}

// AddGroupCommand Adds a sub-command to the root command and configures it using the provided setup function.
func (a *CommandLineApplication) AddGroupCommand(c *cobra.Command, setup types.CommandsSetupFunc) *CommandLineApplication {
	a.root.AddCommand(c)
	if setup != nil {
		wrapper := newCommandSetup(c)
		setup(wrapper)
	}
	return a
}
