package main

import (
	"os"

	"github.com/matzefriedrich/cobra-extensions-example/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
)

func main() {

	app := charmer.NewCommandLineApplication()

	app.
		AddCommand(commands.CreateHelloCommand(), nil).
		AddCommand(commands.CreateCryptCommand(), func(crypto charmer.CommandSetup) {
			crypto.
				AddCommand(commands.CreateEncryptMessageCommand(), nil).
				AddCommand(commands.CreateDecryptMessageCommand(), nil)
		})

	err := app.Execute()
	if err != nil {
		return
	}

	os.Exit(0)
}
