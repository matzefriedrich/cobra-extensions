# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - cobra-extensions v0.3.0, 2024-09-26

### Added

* Introduced the `types` package to centralize command type definitions.

### Changed

* Updated import paths and references across the project to use the new `types` package.
* Removed type definitions from the `commands` and `reflection` package; moved interface types to the new `types` package.
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