package commands

import (
	"github.com/matzefriedrich/cobra-extensions/abstractions"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	abstractions.BaseCommand
	use pkg.CommandName `flag:"crypt"`
}

func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: abstractions.BaseCommand{},
	}
	return pkg.CreateTypedCommand(instance)
}
