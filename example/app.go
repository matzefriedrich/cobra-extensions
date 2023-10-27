package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	app := &cobra.Command{}

	app.AddCommand(
		commands.CreateEncryptMessageCommand(),
		commands.CreateDecryptMessageCommand())

	err := app.Execute()
	if err != nil {
		return
	}

	os.Exit(0)
}
