# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.2.0] - cobra-extensions v0.2.0, 2024-04-25

### Added

- You can now tame the cobra with the `charmer` module; see `src/example/cmd/charmer.go`
- Support for positional arguments (pass values to commands without named flags); see `src/example/cmd/hello.go`

### Changed

- Turns the repository into a multi-module workspace (examples are now separated); run `go work sync` after checkout
- Refactors the whole package (separates types and functions into several new smalled module types), but the package remains compatible with the previous version


## [0.1.0] - cobra-extensions v0.1.0, 2023-10-27

### Added 

- Initial implementation of `cobra-extension` package that provides functionality to define commands and flags in a declarative manner
- Supports command inheritance (define base commands and flags, and reuse them to void duplicate definitions)
- Adds an example application