# depcaps - map capabilities of dependencies agains a set of allowed capabilities

[![Test Status](https://github.com/breml/depcaps/workflows/Go%20Matrix/badge.svg)](https://github.com/breml/depcaps/actions?query=workflow%3AGo%20Matrix) [![Go Report Card](https://goreportcard.com/badge/github.com/breml/depcaps)](https://goreportcard.com/report/github.com/breml/depcaps) [![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

depcaps maps capabilities of dependencies agains a set of allowed capabilities.

## Installation

Download `depcaps` from the [releases](https://github.com/breml/depcaps/releases) or get the latest version from source with:

```shell
go get github.com/breml/depcaps/cmd/depcaps
```

## Usage

### Shell

Check everything:

```shell
depcaps ./...
```

## Inspiration

* [capslock](https://github.com/google/capslock)
* [Capslock: What is your code really capable of?](https://security.googleblog.com/2023/09/capslock-what-is-your-code-really.html)
