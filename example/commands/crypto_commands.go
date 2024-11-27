package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type cryptoCommand struct {
	types.BaseCommand
	use types.CommandName `flag:"crypt" short:"Provides commands to encrypt and decrypt messages." long:"The command group offers tools for securely encrypting messages and decrypting them with a provided passphrase. Use these commands to ensure message confidentiality and safe data transmission."`
}

// CreateCryptCommand initializes a new cobra.Command, providing sub-commands for encrypting and decrypting messages.
func CreateCryptCommand() *cobra.Command {
	instance := &cryptoCommand{
		BaseCommand: types.BaseCommand{},
	}
	return commands.CreateTypedCommand(instance)
}
