# Singlet

[![Build Status](https://img.shields.io/github/actions/workflow/status/jhalterman/singlet/test.yml)](https://github.com/jhalterman/singlet/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/jhalterman/singlet/branch/main/graph/badge.svg?token=VKVY1VTA1U)](https://codecov.io/gh/jhalterman/singlet)
[![License](http://img.shields.io/:license-apache-brightgreen.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Godoc](https://pkg.go.dev/badge/github.com/jhalterman/singlet)](https://pkg.go.dev/github.com/jhalterman/singlet)

Singlet provides threadsafe, generic singletons for Golang.

## Motivation

The common pattern of using `sync.Once` to implement a global singleton doesn't work well when the singleton value includes a generic type. Singlet provides a solution, allowing you to create generic, threadsafe singletons, anywhere in your code.

## Example

To use singlet, first create a `Singleton`. Then call `GetOrDo` which will create and store a value in the `Singleton`, if one doesn't already exist, by calling the provided `func`, else it will return the existing value:

```go
makeFoo := func() *Foo[int] {
    return &Foo[int]{t: 1}
}

var s = &singlet.Singleton{}
foo1, _ := singlet.GetOrDo(s, makeFoo) 
foo2, _ := singlet.GetOrDo(s, makeFoo) // makeFoo is only called once

if foo1 != foo2 {
    panic("foos should be equal")
}
```

You can also get a previously created value for a `Singleton`:

```go
foo, _ := singlet.Get[Foo](singleton)
```

### Possible Errors

Calling `Get` or `GetOrDo` for a result type that doesn't match the previously stored result type for a `Singleton` will result in an `ErrTypeMismatch`:

```go
s := &singlet.Singleton{}
singlet.GetOrDo(s, func() int {
    return 1
})
_, err := singlet.Get[string](s) // Returns ErrTypeMismatch
```

## License

Copyright Jonathan Halterman. Released under the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html).