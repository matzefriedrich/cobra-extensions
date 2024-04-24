package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	pkg.BaseCommand
	use pkg.CommandName `flag:"crypt"`
}

func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: pkg.BaseCommand{},
	}
	return pkg.CreateTypedCommand(instance)
}
