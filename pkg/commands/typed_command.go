package commands

import (
	"context"

	"github.com/matzefriedrich/cobra-extensions/internal/reflection"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type commandContextValue struct {
	handler    types.TypedCommand
	descriptor types.CommandDescriptor
}

type CommandOption func(cmd *cobra.Command)

// run Binds argument and flag values and executes the command.
func (c *commandContextValue) run(ctx context.Context, target *cobra.Command, args ...string) {
	c.descriptor.UnmarshalFlagValues(target)
	c.descriptor.UnmarshalArgumentValues(args...)
	c.handler.Execute(ctx)
}

// CreateTypedCommand Creates a new typed command from the given handler instance.
func CreateTypedCommand[T types.TypedCommand](instance T, options ...func() CommandOption) *cobra.Command {

	reflector := reflection.NewCommandReflector[T]()
	desc := reflector.ReflectCommandDescriptor(instance)

	cmd := &cobra.Command{
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			contextValue := &commandContextValue{
				handler:    instance,
				descriptor: desc,
			}
			ctx := cmd.Context()
			contextValue.run(ctx, cmd, args...)
		},
	}

	for _, option := range options {
		f := option()
		f(cmd)
	}

	desc.BindArguments(cmd)
	desc.BindFlags(cmd)

	return cmd
}

// NonRunnable disables the Run and RunE functions of a Cobra command, effectively making the command non-runnable.
func NonRunnable() CommandOption {
	return func(c *cobra.Command) {
		c.Run = nil
		c.RunE = nil
	}
}
