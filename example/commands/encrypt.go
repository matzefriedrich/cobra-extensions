package commands

import (
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/matzefriedrich/cobra-extensions/v0/pkg"
	"github.com/matzefriedrich/cobra-extensions/v0/pkg/abstractions"
	"github.com/spf13/cobra"
)

var _ abstractions.ExecutableCommand = &encryptMessageCommand{}

type encryptMessageCommand struct {
	cryptCommand
	use     abstractions.CommandName `flag:"encrypt" short:"Encrypt a message." long:"Encrypt a message and protects it with a passphrase."`
	Message string                   `flag:"message" usage:"The message to encrypt."`
}

func CreateEncryptMessageCommand() *cobra.Command {
	instance := &encryptMessageCommand{cryptCommand: cryptCommand{}}
	return pkg.CreateTypedCommand(instance)
}

func (e *encryptMessageCommand) Execute() {

	message := crypto.NewPlainMessageFromString(e.Message)
	encrypted, _ := crypto.EncryptMessageWithPassword(message, []byte(e.Passphrase))
	armored, _ := encrypted.GetArmored()
	_, _ = os.Stdout.WriteString(armored)
}
