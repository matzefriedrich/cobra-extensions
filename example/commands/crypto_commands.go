package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	types.BaseCommand
	use types.CommandName `flag:"crypt"`
}

func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: types.BaseCommand{},
	}
	return commands.CreateTypedCommand(instance)
}
