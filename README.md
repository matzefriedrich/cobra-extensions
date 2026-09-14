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


## Features

- **Struct-based commands** - define commands and flags as annotated Go struct fields instead of writing manual boilerplate.
- **Typed command handlers** - implement a simple `Execute(ctx)` method instead of wiring up `cobra.Command.Run` by hand.
- **Command grouping** - organize sub-commands hierarchically via a declarative setup function.
- **Command inheritance** - define base commands and flags once, and reuse them across commands to avoid duplicate definitions.
- **Positional arguments** - bind unnamed command-line values to struct fields alongside regular flags.
- **Slice-valued flags** - map flags specified multiple times on the command line to slice fields (`[]string`, `[]int`, `[]int64`, `[]bool`).
- **Markdown documentation generation** - generate Markdown docs for all registered commands and subcommands via `NewMarkdownCommand`.
- **Default value providers** - let flags fall back to values from environment variables, YAML files, or any custom source (e.g., a vault backend) when not set on the command line.
- **Pluggable cobra-x tag parsers** - choose between the compact DSL and conventional space-separated struct tags, or plug in a custom `types.CobraXTagParser`.


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
    types.BaseCommand `cobra-x:"hello, help='Prints a greeting.'"`
    Name              string `cobra-x:"-n|--name, help='The name to greet', default='World'"`
}

func (g *greetCommand) Execute(_ context.Context) {
    fmt.Printf("Hello, %s!\n", g.Name)
}

func NewGreetCommand() *cobra.Command {
    return commands.CreateTypedCommand(&greetCommand{})
}
```

Can you spot the difference? There are no global vars, no manual flag binding, no extra ceremony. This pattern improves clarity, testability, and maintainability - especially as your CLI grows beyond a handful of commands.

> [!IMPORTANT]
> Since version `v0.7.0`, `types.CommandName` is deprecated. Use `types.BaseCommand` and the `cobra-x` tag for defining command metadata instead.

## How struct tags become CLI flags

Each field in your command struct becomes a flag. You control its name, description, and shorthand through the `cobra-x` struct tag:

```go
Name string `cobra-x:"-n|--name, help='The name to greet', default='World'"`
```

The `cobra-x` tag supports:
- Simplified name expressions: `-n|--name`, `--name`, or just `-n`.
- Key-value attributes: `help`, `description`, `default`.
- Backward compatibility with legacy tags like `name`, `shorthand`, `usage`, and `default`.

See [https://github.com/matzefriedrich/cobra-extensions-docs](https://github.com/matzefriedrich/cobra-extensions-docs) for complete usage examples.


## Pluggable cobra-x tag parsers

The `cobra-x` tags can be read in two interchangeable syntaxes: the compact DSL shown above and conventional space-separated struct tags. Both are implementations of the `types.CobraXTagParser` interface, so you can even plug in a parser of your own. The compact syntax keeps tags short and scannable; the conventional syntax uses fixed, well-known tag keys that can be validated at compile time. Pick the route that fits your project - the compact DSL is the default.

To use the conventional form for all commands of an application, inject `types.NewStandardTagParser()` before adding typed commands:

```go
app := charmer.NewCommandLineApplication("demo", "A short description.")
app.WithTagParser(types.NewStandardTagParser())
app.AddTypedCommand(&greetCommand{})
```

Alternatively, give a single command its own parser:

```go
commands.CreateTypedCommandWithTagParser(&greetCommand{}, types.NewStandardTagParser())
```

> [!NOTE]
> Commands are reflected at creation time, so commands created through an application must be added via `AddTypedCommand` for the application's parser to take effect.

See [https://github.com/matzefriedrich/cobra-extensions-docs](https://github.com/matzefriedrich/cobra-extensions-docs) for the conventional tag syntax and complete usage examples.


## Positional arguments

Besides flags, you can bind unnamed command-line values to struct fields. Declare a plain struct that embeds a `types.CommandArgs` field, followed by the positional fields in the order they appear on the command line:

```go
type greetArgs struct {
    types.CommandArgs
    Name string `cobra-x:"name"`
}

type greetCommand struct {
    types.BaseCommand `cobra-x:"greet, help='Greet someone'"`
    Arguments         greetArgs
}
```

Each positional field is bound in declaration order, so `greet Alice` sets `Arguments.Name` to `"Alice"`.

To require a minimum number of positional arguments, configure the `CommandArgs` field:

```go
Arguments greetArgs

// ...

Arguments: greetArgs{CommandArgs: types.NewCommandArgs(types.MinimumArgumentsRequired(1))},
```

Since version `v0.10.0`, a positional field with a `cobra-x` tag (e.g. `cobra-x:"name"`) renders its placeholder in the command's `Use` string and the `--help` Usage line. Placeholders in required positions render as `<name>`; the remaining ones render as `[name]`:

```
Usage:
  demo greet <name> [flags]
```

Untagged positional fields still bind at runtime but are omitted from the `Use` string.

See [https://github.com/matzefriedrich/cobra-extensions-docs](https://github.com/matzefriedrich/cobra-extensions-docs) for complete usage examples.


## Default value providers

Since version `v0.8.0`, flags can automatically receive their values from an external configuration source (environment variables, YAML files, or a custom backend like a vault service) whenever they are not explicitly set on the command line. A flag opts in via the `setting-key` attribute on its `cobra-x` tag, and a `DefaultValueProvider` is injected into the application to resolve values, with the command line always taking precedence.

See [https://github.com/matzefriedrich/cobra-extensions-docs](https://github.com/matzefriedrich/cobra-extensions-docs) for the full list of built-in providers and usage examples.


## Design Philosophy

Cobra gives you flexibility - and with it, a lot of responsibility. You're on your own when it comes to structuring CLIs, managing shared configuration, or avoiding boilerplate; the `cobra-extensions` package takes a more opinionated stance:

- Commands should be self-contained structs.

- Flags should be defined where they're used - at the command level, not in a separate init method somewhere else in the code.

- Manual flag wiring is a waste of time, and does not guide with consistency, which is more valuable than flexibility for most use cases.

- Struct-based design makes testing and composition easier. For instance, flag definitions can be reused via embeddings, avoiding the redefinition of flags per command.

- Step down to raw Cobra at any time, `cobra-extensions` won´t prevent you from doing so.


## Development

```bash
make install-tools   # install golangci-lint and goveralls into ./bin
make lint            # run the linter
make coverage        # run the CI coverage command and print the total
```

`make coverage` mirrors the CI pipeline exactly (cross-package `-coverpkg ./...` profile), so you can check coverage locally the same way the badge is computed.

---

Copyright 2023 - 2026 by Matthias Friedrich