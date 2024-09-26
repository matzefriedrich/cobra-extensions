[![CI](https://github.com/matzefriedrich/cobra-extensions/actions/workflows/go.yml/badge.svg)](https://github.com/matzefriedrich/cobra-extensions/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/matzefriedrich/cobra-extensions/badge.svg)](https://coveralls.io/github/matzefriedrich/cobra-extensions)
[![Go Reference](https://pkg.go.dev/badge/github.com/matzefriedrich/cobra-extensions.svg)](https://pkg.go.dev/github.com/matzefriedrich/cobra-extensions)
[![Go Report Card](https://goreportcard.com/badge/github.com/matzefriedrich/cobra-extensions)](https://goreportcard.com/report/github.com/matzefriedrich/cobra-extensions)
![License](https://img.shields.io/github/license/matzefriedrich/cobra-extensions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matzefriedrich/cobra-extensions)
![GitHub Release](https://img.shields.io/github/v/release/matzefriedrich/cobra-extensions?include_prereleases)

# cobra-extensions

**cobra-extensions** is an extension package for the well-known [Cobra](https://github.com/spf13/cobra) library, designed to enhance command management by introducing a declarative approach for binding flags to command structs. With this approach, you build complex CLI applications consisting of many commands in a clean and organized manner.

By leveraging command structs, adopting patterns like dependency injection becomes far easier, further simplifying the design of scalable and maintainable CLI tools. The library automates flag setup using struct tags, reducing boilerplate and improving developer productivity.


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