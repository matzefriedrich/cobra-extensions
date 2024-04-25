package commands

import (
	"github.com/matzefriedrich/cobra-extensions/v0/pkg"
	"github.com/matzefriedrich/cobra-extensions/v0/pkg/abstractions"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	abstractions.BaseCommand
	use abstractions.CommandName `flag:"crypt"`
}

func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: abstractions.BaseCommand{},
	}
	return pkg.CreateTypedCommand(instance)
}
