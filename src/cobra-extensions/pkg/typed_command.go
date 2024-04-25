package pkg

import (
	"context"
	"github.com/matzefriedrich/cobra-extensions/pkg/reflection"
	"github.com/spf13/cobra"
)

type TypedCommand interface {
	Execute()
}

type commandContextValue struct {
	handler    TypedCommand
	descriptor reflection.CommandDescriptor
}

// run Binds argument and flag values and executes the command.
func (c *commandContextValue) run(target *cobra.Command, args ...string) {
	c.descriptor.UnmarshalFlagValues(target)
	c.descriptor.UnmarshalArgumentValues(args...)
	c.handler.Execute()
}

// CreateTypedCommand Creates a new typed command from the given handler instance.
func CreateTypedCommand[T TypedCommand](instance T) *cobra.Command {

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
