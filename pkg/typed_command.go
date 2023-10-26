package pkg

import (
	"context"
	"github.com/spf13/cobra"
)

type TypedCommand interface {
	Execute()
}

// CreateTypedCommand Creates a new typed command from the given handler instance.
func CreateTypedCommand[T TypedCommand](instance T) *cobra.Command {

	desc := ReflectCommandDescriptor(instance)

	cmd := &cobra.Command{
		Use: desc.use,
		Run: func(cmd *cobra.Command, args []string) {
			handler := cmd.Context().Value(desc.key).(T)
			UnmarshalCommand(cmd, desc)
			handler.Execute()
		},
	}

	desc.BindFlags(cmd)

	ctx := context.WithValue(context.Background(), desc.key, instance)
	cmd.SetContext(ctx)

	return cmd
}
