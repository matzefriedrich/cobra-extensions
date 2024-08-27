package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"log"
)

func main() {

	err :=
		charmer.NewCommandLineApplication("charmer-example", "").
			AddCommand(commands.CreateHelloCommand()).
			AddGroupCommand(commands.CreateCryptCommand(), func(crypto charmer.CommandSetup) {
				crypto.AddCommand(
					commands.CreateEncryptMessageCommand(),
					commands.CreateDecryptMessageCommand())
			}).
			Execute()

	if err != nil {
		log.Fatal(err)
	}
}
