package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	app := &cobra.Command{}

	app.AddCommand(commands.CreateHelloCommand())

	AddCryptCommands(app)

	err := app.Execute()
	if err != nil {
		return
	}

	os.Exit(0)
}

func AddCryptCommands(app *cobra.Command) {

	crypt := commands.CreateCryptCommand()

	crypt.AddCommand(
		commands.CreateEncryptMessageCommand(),
		commands.CreateDecryptMessageCommand(),
	)

	app.AddCommand(crypt)
}
