package pkg

import (
	"context"
	"github.com/matzefriedrich/cobra-extensions/pkg/reflection"
	"github.com/spf13/cobra"
)

type TypedCommand interface {
	Execute()
}

// CreateTypedCommand Creates a new typed command from the given handler instance.
func CreateTypedCommand[T TypedCommand](instance T) *cobra.Command {

	reflector := reflection.NewCommandReflector[T]()
	desc := reflector.ReflectCommandDescriptor(instance)

	commandKey := desc.Key()

	cmd := &cobra.Command{
		Use:   desc.Use(),
		Short: desc.ShortDescriptionText(),
		Long:  desc.LongDescriptionText(),
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(commandKey).(T)
			reflection.UnmarshalCommand(cmd, desc, args...)
			handler.Execute()
		},
	}

	cmd.Args = cobra.MinimumNArgs(desc.Arguments().MinimumArgs)

	desc.BindFlags(cmd)

	ctx := context.WithValue(context.Background(), commandKey, instance)
	cmd.SetContext(ctx)

	return cmd
}
