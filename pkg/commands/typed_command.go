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

// run Binds argument and flag values and executes the command.
func (c *commandContextValue) run(target *cobra.Command, args ...string) {
	c.descriptor.UnmarshalFlagValues(target)
	c.descriptor.UnmarshalArgumentValues(args...)
	c.handler.Execute()
}

// CreateTypedCommand Creates a new typed command from the given handler instance.
func CreateTypedCommand[T types.TypedCommand](instance T) *cobra.Command {

	reflector := reflection.NewCommandReflector[T]()
	desc := reflector.ReflectCommandDescriptor(instance)

	commandKey := desc.Key()

	cmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			v := cmd.Context().Value(commandKey).(*commandContextValue)
			v.run(cmd, args...)
		},
	}

	desc.BindArguments(cmd)
	desc.BindFlags(cmd)

	contextValue := &commandContextValue{
		handler:    instance,
		descriptor: desc,
	}

	ctx := context.WithValue(context.Background(), commandKey, contextValue)
	cmd.SetContext(ctx)

	return cmd
}
