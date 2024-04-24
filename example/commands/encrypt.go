package commands

import (
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
	"os"
)

var _ pkg.ExecutableCommand = &encryptMessageCommand{}

type encryptMessageCommand struct {
	cryptCommand
	use     pkg.CommandName `flag:"encrypt" short:"Encrypt a message." long:"Encrypt a message and protects it with a passphrase."`
	Message string          `flag:"message" usage:"The message to encrypt."`
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
