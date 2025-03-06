# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.5] - cobra-extensions v0.4.5, 2025-03-06

### Changed

* Bumps github.com/ProtonMail/gopenpgp/v2 from 2.8.2 to 2.8.3 [#17](https://github.com/matzefriedrich/cobra-extensions/pull/17)


## [0.4.4] - cobra-extensions v0.4.4, 2025-02-19

### Changed

* Bumps github.com/spf13/cobra from 1.8.1 to 1.9.1 [#16](https://github.com/matzefriedrich/cobra-extensions/pull/16)


## [0.4.3] - cobra-extensions v0.4.3, 2025-01-13

### Changed

* Bumps `github.com/ProtonMail/gopenpgp/v2` from 2.8.1 to 2.8.2 [#15](https://github.com/matzefriedrich/cobra-extensions/pull/15)


## [0.4.2] - cobra-extensions v0.4.2, 2024-12-16

### Changed

* Bumps `golang.org/x/crypto` from 0.17.0 to 0.31.0 [#14](https://github.com/matzefriedrich/cobra-extensions/pull/14)


## [0.4.1] - cobra-extensions v0.4.1, 2024-12-03

### Changed

* Bumps `github.com/ProtonMail/gopenpgp/v2` from 2.8.0 to 2.8.1 [#13](https://github.com/matzefriedrich/cobra-extensions/pull/13)


## [0.4.0] - cobra-extensions v0.4.0, 2024-11-27

### Added

* Introduced the `NewMarkdownCommand` method to create a `cobra.Command` instance for generating Markdown documentation. This command can be linked to the root command and produces documentation for all registered commands and subcommands in Markdown format. [#12](https://github.com/matzefriedrich/cobra-extensions/pull/12)

### Changed

* Enhanced the `CreateTypedCommand` function to support options, enabling greater flexibility in command creation and configuration. This update allows group commands to be marked as non-runnable, influencing their representation in Markdown documentation. To achieve this, pass the `NonRunnable` option to the `CreateTypedCommand` function.

* Improved and expanded descriptions for example command groups and applications to enhance clarity and usability.


## [0.3.3] - cobra-extensions v0.3.3, 2024-11-27

### Changed

* Bumps `github.com/stretchr/testify` from 1.9.0 to 1.10.0 [#11](https://github.com/matzefriedrich/cobra-extensions/pull/11)


## [0.3.2] - cobra-extensions v0.3.2, 2024-11-12

### Changed

* Bumps `github.com/ProtonMail/gopenpgp/v2` from 2.7.5 to 2.8.0 [#10](https://github.com/matzefriedrich/cobra-extensions/pull/10)


## [0.3.1] - cobra-extensions v0.3.1, 2024-09-26

### Added

* Added tests and increases coverage [#9](https://github.com/matzefriedrich/cobra-extensions/pull/9)

### Changed

* Updated dependencies (as suggested by dependabot; see PR #6, #7, and #8) to address detected security issues
* Minor refactorings to the `command_setup.go` module


## [0.3.0] - cobra-extensions v0.3.0, 2024-09-26

### Added

* Introduced the `types` package to centralize command type definitions.

### Changed

* Updated import paths and references across the project to use the new `types` package.
* Removed type definitions from the `commands` and `reflection` packages; moved interface types to the new `types` package.
* The `reflection` package is now internal because it supports the functionality of the extensions package and is not intended for use by user code.
* Upgraded Go SDK version from 1.21 to 1.23.


## [0.2.6] - cobra-extensions v0.2.6, 2024-08-27

### Added

- Adds the `NewRootCommand` helper method that can be used to create root `cobra.Command` objects with application name and description.

### Changed

- Changes the `NewCommandLineApplication` method to set the application name and description (required to generate proper completion scripts for apps)

- Adopts changes to examples


## [0.2.5] - cobra-extensions v0.2.5, 2024-06-12

### Fixes

- The `reflection` module ignores fields of unsupported types


## [0.2.3] - cobra-extensions v0.2.3, 2024-04-25

### Added

- You can now tame the cobra with the `charmer` module; see `example/cmd/charmer/charmer.go`
- Support for positional arguments (pass values to commands without named flags); see `example/commands/hello.go`

### Changed

- Turns the repository into a multi-module workspace (examples are now separated); run `go work sync` after checkout
- Refactors the whole package (separates types and functions into several new smalled module types), but the package remains compatible with the previous version


## [0.1.0] - cobra-extensions v0.1.0, 2023-10-27

### Added 

- Initial implementation of `cobra-extension` package that provides functionality to define commands and flags in a declarative manner
- Supports command inheritance (define base commands and flags, and reuse them to void duplicate definitions)
- Adds an example application