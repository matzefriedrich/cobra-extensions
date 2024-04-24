package main

import (
	"os"

	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/spf13/cobra"
)

func main() {

	app := NewCommandLineApplication()

	app.
		AddCommand(commands.CreateHelloCommand(), nil).
		AddCommand(commands.CreateCryptCommand(), func(crypto CommandSetup) {
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
