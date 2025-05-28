[![CI](https://github.com/matzefriedrich/cobra-extensions/actions/workflows/go.yml/badge.svg)](https://github.com/matzefriedrich/cobra-extensions/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/matzefriedrich/cobra-extensions/badge.svg)](https://coveralls.io/github/matzefriedrich/cobra-extensions)
[![Go Reference](https://pkg.go.dev/badge/github.com/matzefriedrich/cobra-extensions.svg)](https://pkg.go.dev/github.com/matzefriedrich/cobra-extensions)
[![Go Report Card](https://goreportcard.com/badge/github.com/matzefriedrich/cobra-extensions)](https://goreportcard.com/report/github.com/matzefriedrich/cobra-extensions)
![License](https://img.shields.io/github/license/matzefriedrich/cobra-extensions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/matzefriedrich/cobra-extensions)
![GitHub Release](https://img.shields.io/github/v/release/matzefriedrich/cobra-extensions?include_prereleases)


# cobra-extensions

**cobra-extensions** is an opinionated extension library for [spf13/cobra](https://github.com/spf13/cobra) that takes a more structured and declarative approach to building CLI applications. It simplifies defining commands and flags by letting you use Go structs with annotated fields, avoiding repetitive flag setup and glue code.

Suppose you've built large CLI tools with Cobra. In that case, you've likely hit its boilerplate-heavy design, manual flag binding, and lack of structure around command logic. `cobra-extensions` addresses these pain points by binding flags directly to fields via struct tags, so commands become just regular Go structs with minimal ceremony.

This approach encourages clean separation of concerns, easier testability, and better scaling as the number of commands grows. The trade-off? You give up some of Cobra’s flexibility for a streamlined, consistent pattern that works well for most real-world use cases.

If you're okay with that, this package will save you time, and your CLI codebase will thank you.


## Usage

Using the Cobra extensions is a no-brainer. Use, `go get` to install the latest version of the library.

````bash
go get -u github.com/matzefriedrich/cobra-extensions@latest
````


## Why cobra-extensions?

Cobra is mighty, but as your CLI grows, its unstructured nature starts to show: flags are managed through global variables, commands require repetitive boilerplate, and there's little guidance on organizing logic cleanly.  `cobra-extensions` introduces a declarative, opinionated approach that eliminates manual flag wiring and encourages structured command definitions. Let's compare the differences:

### Traditional cobra

```go
var name string

var greetCmd = &cobra.Command{
    Use:   "hello",
    Short: "Prints a greeting.",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("Hello, %s!\n", name)
    },
}

func init() {
    greetCmd.Flags().StringVarP(&name, "name", "n", "World", "The name to greet")
}
```

### With cobra-extensions

```go
type greetCommand struct {
    use  types.CommandName `flag:"hello" short:"Prints a greeting."`
    Name string            `flag:"name" short:"The name to greet." default:"World"`
}

func (g *greetCommand) Execute() {
    fmt.Printf("Hello, %s!\n", g.Name)
}

func NewGreetCommand() *cobra.Command {
    return commands.CreateTypedCommand(&greetCommand{})
}
```

Can you spot the difference? There are no global vars, no manual flag binding, no extra ceremony. This pattern improves clarity, testability, and maintainability - especially as your CLI grows beyond a handful of commands.

## How struct tags become CLI flags

Each field in your command struct becomes a flag. You control its name, description, and default value through struct tags:

```go
Name string `flag:"name" short:"The name to greet." default:"World"`
```

See [https://github.com/matzefriedrich/cobra-extensions-docs](https://github.com/matzefriedrich/cobra-extensions-docs) for complete usage examples.


## Design Philosophy

Cobra gives you flexibility - and with it, a lot of responsibility. You're on your own when it comes to structuring CLIs, managing shared configuration, or avoiding boilerplate; the `cobra-extensions` package takes a more opinionated stance:

- Commands should be self-contained structs.

- Flags should be defined where they're used - at the command level, not in a separate init method somewhere else in the code.

- Manual flag wiring is a waste of time, and does not guide with consistency, which is more valuable than flexibility for most use cases.

- Struct-based design makes testing and composition easier. For instance, flag definitions can be reused via embeddings, avoiding the redefinition of flags per command.

- Step down to raw Cobra at any time, `cobra-extensions` won´t prevent you from doing so.

---

Copyright 2023 - 2025 by Matthias Friedrich