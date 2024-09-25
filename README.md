# cobra-extensions

**An extension library for Cobra**, adding a declarative approach to simplify the configuration of commands and flags.

## Usage

Using the Cobra extensions is a no-brainer. Use, `go get` to install the latest version of the library.

````bash
$ go get -u github.com/matzefriedrich/cobra-extensions@latest
````

Next, include Cobra extensions in your application.

````go
import (
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)
````

## Example

All it needs is the definition of a struct that implements the `TypedCommand` interface. Public fields are annotated with `flag` tags to define flags and bind them to the struct once parsed. A private field of type `CommandName` specifies the command's name. For instance:

````golang
package commands

import (
	"fmt"
	"github.com/matzefriedrich/cobra-extensions/pkg/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/types"
	"github.com/spf13/cobra"
)

type helloCommand struct {
	use       types.CommandName `flag:"hello"`
	Arguments helloArgs
}

var _ types.TypedCommand = (*helloCommand)(nil)

type helloArgs struct {
	types.CommandArgs
	Name string
}

func CreateHelloCommand() *cobra.Command {
	instance := &helloCommand{
		Arguments: helloArgs{
			CommandArgs: types.NewCommandArgs(types.MinimumArgumentsRequired(1)),
		}}
	return commands.CreateTypedCommand(instance)
}

func (c *helloCommand) Execute() {
	fmt.Printf("Hello %s.\n", c.Arguments.Name)
}

````

A `CreateHelloCommand` factory method creates a new `helloCommand` instance and utilizes the `CreateTypedCommand` method to create and initialize a Cobra command.

````golang
package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/matzefriedrich/cobra-extensions/pkg/charmer"
	"github.com/spf13/cobra"
)

func main() {

	app := charmer.NewRootCommand("simple-example", "")

	app.AddCommand(commands.CreateHelloCommand())

	err := app.Execute()
	if err != nil {
		panic(err)
	}
}
````