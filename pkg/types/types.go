package types

import (
	"context"

	"github.com/spf13/cobra"
)

type CommandDescriptor interface {
	UnmarshalFlagValues(target *cobra.Command)
	UnmarshalArgumentValues(args ...string)
	BindArguments(cmd *cobra.Command)
	BindFlags(cmd *cobra.Command)

	// SetDefaultValueProvider sets the provider used to retrieve default values for command flags.
	SetDefaultValueProvider(provider DefaultValueProvider)
}

// CommandReflector Reflects metadata of a command handler to create a CommandDescriptor.
type CommandReflector[T any] interface {
	ReflectCommandDescriptor(n T) CommandDescriptor
}

// ArgumentsDescriptorOption defines an option that can configure an ArgumentsDescriptor.
type ArgumentsDescriptorOption func(argumentsDescriptor any)

// ArgumentsDescriptor provides an interface for managing command arguments metadata.
type ArgumentsDescriptor interface {
	With(options ...ArgumentsDescriptorOption) ArgumentsDescriptor
	BindArguments(target *cobra.Command)
	BindArgumentValues(args ...string)
}

// TypedCommand represents a typed command that can be executed.
type TypedCommand interface {

	// Execute runs the command, performing its specific action based on the implementation.
	Execute(ctx context.Context)
}

// CommandsSetupFunc defines a function type used to set up commands within a CommandSetup context.
type CommandsSetupFunc func(setup CommandSetup)

// CommandSetup provides methods to add and organize commands within a command-line application interface.
type CommandSetup interface {

	// AddCommand adds one or more sub-commands to the current command.
	AddCommand(c ...*cobra.Command) CommandSetup

	// AddGroupCommand adds a sub-command to the current command and calls the setup function for additional configuration.
	AddGroupCommand(c *cobra.Command, setup CommandsSetupFunc) CommandSetup
}

// DefaultValueProvider provides a mechanism to retrieve default values for command flags from an external source.
type DefaultValueProvider interface {

	// GetValue retrieves the value associated with the specified key.
	// It returns the value if found, or an empty string if not found.
	// An error is returned if the retrieval process fails.
	GetValue(key string) (string, error)
}

type contextKey string

const (
	DefaultValueProviderKey contextKey = "cobra-x-default-value-provider"
)
