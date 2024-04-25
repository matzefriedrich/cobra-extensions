package commands

import (
	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/matzefriedrich/cobra-extensions/pkg/abstractions"
	"io"
	"os"

	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/spf13/cobra"
)

var _ abstractions.ExecutableCommand = &decryptMessageCommand{}

type decryptMessageCommand struct {
	cryptCommand
	use abstractions.CommandName `flag:"decrypt" short:"Decrypt a message." long:"Reads an armored message from stdin and decrypts it."`
}

func CreateDecryptMessageCommand() *cobra.Command {
	instance := &decryptMessageCommand{cryptCommand: cryptCommand{}}
	return pkg.CreateTypedCommand(instance)
}

func (d *decryptMessageCommand) Execute() {
	armored, _ := ReadArmoredMessagedFromStdin()
	message, _ := crypto.NewPGPMessageFromArmored(armored)
	decrypted, _ := crypto.DecryptMessageWithPassword(message, []byte(d.Passphrase))
	_, _ = os.Stdout.WriteString(decrypted.GetString())
}

func ReadArmoredMessagedFromStdin() (string, error) {
	data, err := io.ReadAll(os.Stdin)
	if err == nil {
		return string(data), nil
	}
	return "", err
}
