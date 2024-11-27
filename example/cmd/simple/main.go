package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	builtin "github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	app := charmer.NewRootCommand("simple-example", "")

	app.AddCommand(builtin.NewMarkdownDocsCommand(app))
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
