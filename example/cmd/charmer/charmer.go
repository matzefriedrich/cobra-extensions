package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"log"
)

func main() {

	err :=
		charmer.NewCommandLineApplication("charmer-example", "").
			AddCommand(commands.CreateHelloCommand()).
			AddGroupCommand(commands.CreateCryptCommand(), func(crypto types.CommandSetup) {
				crypto.AddCommand(
					commands.CreateEncryptMessageCommand(),
					commands.CreateDecryptMessageCommand())
			}).
			Execute()

	if err != nil {
		log.Fatal(err)
	}
}
