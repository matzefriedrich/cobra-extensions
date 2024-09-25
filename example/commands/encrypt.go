package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/reflection"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

var _ types.TypedCommand = &encryptMessageCommand{}

type encryptMessageCommand struct {
	cryptCommand
	use     reflection.CommandName `flag:"encrypt" short:"Encrypt a message." long:"Encrypt a message and protects it with a passphrase."`
	Message string                 `flag:"message" usage:"The message to encrypt."`
}

func CreateEncryptMessageCommand() *cobra.Command {
	instance := &encryptMessageCommand{cryptCommand: cryptCommand{}}
	return commands.CreateTypedCommand(instance)
}

func (e *encryptMessageCommand) Execute() {

	message := crypto.NewPlainMessageFromString(e.Message)
	encrypted, _ := crypto.EncryptMessageWithPassword(message, []byte(e.Passphrase))
	armored, _ := encrypted.GetArmored()
	_, _ = os.Stdout.WriteString(armored)
}
