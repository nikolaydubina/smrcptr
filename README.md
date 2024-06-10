## Same Receiver Pointer (smrcptr)

[![go-recipes](https://raw.githubusercontent.com/nikolaydubina/go-recipes/main/badge.svg?raw=true)](https://github.com/nikolaydubina/go-recipes)
[![Go Report Card](https://goreportcard.com/badge/github.com/nikolaydubina/smrcptr)](https://goreportcard.com/report/github.com/nikolaydubina/smrcptr)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/nikolaydubina/smrcptr/badge)](https://securityscorecards.dev/viewer/?uri=github.com/nikolaydubina/smrcptr)

> Don't mix receiver types. Choose either pointers or struct types for all available methods.

Go has rules on how it automatically selects either value or method receivers, which is complex and can lead to bugs.
Therefore, it is a common style recommendation[^1][^2].

```bash
go install github.com/nikolaydubina/smrcptr@latest
```

```go
type Pancake struct{}

func NewPancake() Pancake { return Pancake{} }

func (s *Pancake) Fry() {}

func (s Pancake) Bake() {}
```

```bash
$ smrcptr ./...
/pancake.go:12:1: Pancake.Fry uses pointer
/pancake.go:10:1: Pancake.NewPancake uses value
/pancake.go:14:1: Pancake.Bake uses value
```

## Existing Linters

#### staticcheck

As of `2022-11-30`, it does not detect that pointer and value method receivers are mixed.
Most relevant analyzer `ST1016` checks only name of method receiver.

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
* https://github.com/dominikh/go-tools/issues/911

[^1]: Go wiki https://github.com/golang/go/wiki/CodeReviewComments#receiver-type
[^2]: Google Go style guide https://google.github.io/styleguide/go/decisions#receiver-type
