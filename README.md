## 🥞 Same Receiver Pointer (smrcptr)

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/smrcptr)](https://goreportcard.com/report/github.com/nikolaydubina/smrcptr)

This `go vet` compatible linter detects mixing pointer and value method receivers for the same type.

```go
type Pancake struct{}

func NewPancake() Pancake { return Pancake{} }

func (s *Pancake) Fry() {}

func (s Pancake) Bake() {}
```

```bash
$ smrcptr ./...
smrcptr/internal/bakery/pancake.go:7:1: Pancake.Fry uses pointer
smrcptr/internal/bakery/pancake.go:9:1: Pancake.Bake uses value
```

Why this is useful? Go has rules on how it can automatically select value and method receivers, which is complex and can lead to bugs.
It is also common style recommendation [Go wiki](https://github.com/golang/go/wiki/CodeReviewComments#receiver-type) and [Google Go style guide](https://google.github.io/styleguide/go/decisions#receiver-type):

> Don't mix receiver types. Choose either pointers or struct types for all available methods.

## Requirements

```bash
go install github.com/nikolaydubina/smrcptr@latest
```

## Features

### Return Status Code

When issue is detected, related info is printed and status code is non-zero.
This is similar as other `go vet` and linters.
This allows convenient use in CI.

### Constructor

It is also useful to detect if "construtor" functions that commonly start with `New...` returns value that matches used in receivers.

```bash
$ smrcptr --constructor=true ./internal/...
smrcptr/internal/bakery/pancake.go:7:1: Pancake.Fry uses pointer
smrcptr/internal/bakery/pancake.go:5:1: Pancake.NewPancake uses value
smrcptr/internal/bakery/pancake.go:9:1: Pancake.Bake uses value
smrcptr/internal/bakery/pancake.go:14:1: Cake.Fry uses pointer
smrcptr/internal/bakery/pancake.go:16:1: Cake.Bake uses value
smrcptr/internal/bakery/pancake.go:23:1: Brownie.Bake uses pointer
smrcptr/internal/bakery/pancake.go:21:1: Brownie.NewBrownie uses value
smrcptr/internal/bakery/pancake.go:35:1: BadCookie.NewBadCookie uses pointer
smrcptr/internal/bakery/pancake.go:37:1: BadCookie.Bake uses value
```

## Existing Linters

#### staticcheck

As of `2022-11-30`, it does not detect that pointer and value method receivers are mixed.
Most relevant analyzser `ST1016` checks only name of method reciver.

```bash
$ staticcheck -checks ST1016 ./...    
main.go:9:18: methods on the same type should have the same receiver name (seen 1x "v", 2x "s") (ST1016)
```

Using all analyzers does not detect it either.

```bash
staticcheck -checks all ./...
main.go:9:18: methods on the same type should have the same receiver name (seen 1x "v", 2x "s") (ST1016)
```

## References

* https://github.com/golang/go/wiki/CodeReviewComments#receiver-type
* https://golang.org/ref/spec#Method_declarations
* https://golangci-lint.run/usage/linters/
* https://github.com/dominikh/go-tools/blob/master/stylecheck/lint.go#L295
* https://github.com/dominikh/go-tools/tree/master/stylecheck/testdata/src/CheckReceiverNamesIdentical
* https://google.github.io/styleguide/go/decisions#receiver-type
