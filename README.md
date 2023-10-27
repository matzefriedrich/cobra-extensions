# cobra-extensions

**An extension library for Cobra**, adding a declarative approach to simplify the configuration of commands and flags.

## Usage

Using the Cobra extensions is a no-brainer. Use, `go get` to install the version of the library.

````bash
$ go get -u github.com/matzefriedrich/cobra-extensions@latest
````

Next, include Cobra extensions in your application.

````go
import "github.com/matzefriedrich/cobra-extensions/pkg"
````

## Example

All it needs is the definition of a struct that implements the `TypedCommand` interface. Public fields are annotated with `flag` tags to define flags and bind them to the struct once parsed. A private field of type `CommandName` specifies the command's name. For instance:

````golang
package commands

import (
	"fmt"

	"github.com/matzefriedrich/cobra-extensions/pkg"
	"github.com/spf13/cobra"
)

type helloCommand struct {
	use  pkg.CommandName `flag:"hello"`
	Name string          `flag:"name" usage:"Your name"`
}

func CreateHelloCommand() *cobra.Command {
	instance := &helloCommand{}
	return pkg.CreateTypedCommand(instance)
}

func (c *helloCommand) Execute() {
	_ = fmt.Sprintf("Hello %s.", c.Name)
}
````

A `CreateHelloCommand` factory method creates a new `helloCommand` instance and utilizes the `CreateTypedCommand` method to create and initialize a Cobra command.

````golang
package main

import (
	"github.com/matzefriedrich/cobra-extensions/example/commands"
	"github.com/spf13/cobra"
	"os"
)

func main() {

    app := &cobra.Command{}

    app.AddCommand(commands.CreateHelloCommand())
	
    _ = app.Execute()
}
````