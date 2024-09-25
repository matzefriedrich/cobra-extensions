package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/reflection"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	commands.BaseCommand
	use reflection.CommandName `flag:"crypt"`
}

func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: commands.BaseCommand{},
	}
	return commands.CreateTypedCommand(instance)
}
