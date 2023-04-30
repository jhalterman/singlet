# Singlet

[![Build Status](https://img.shields.io/github/actions/workflow/status/jhalterman/singlet/test.yml)](https://github.com/jhalterman/singlet/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/jhalterman/singlet/branch/main/graph/badge.svg?token=VKVY1VTA1U)](https://codecov.io/gh/jhalterman/singlet)
[![License](http://img.shields.io/:license-apache-brightgreen.svg)](http://www.apache.org/licenses/LICENSE-2.0.html)
[![Godoc](https://pkg.go.dev/badge/github.com/jhalterman/singlet)](https://pkg.go.dev/github.com/jhalterman/singlet)

Singlet provides threadsafe, generic singletons for Golang.

## Motivation

The common pattern of using `sync.Once` to implement a global singleton doesn't work well when the singleton value includes a generic type, since type arguments may not be available at a global level. Singlet provides a solution, allowing you to create threadsafe, generic singletons, anywhere in your code.

## Example

To use singlet, first create a `Singleton`. Then call `GetOrDo` which will create and store a value in the `Singleton`, if one doesn't already exist, by calling the provided `func`, else it will return the existing value:

```go
var s = &singlet.Singleton{}
cache1, _ := singlet.GetOrDo(s, cache.New[int]) 
cache2, _ := singlet.GetOrDo(s, cache.New[int]) // makeFoo is only called once

if cache1 != cache2 {
    panic("caches should be equal")
}
```

You can also get a previously created value for a `Singleton`:

```go
cache, _ := singlet.Get[*Cache[int]](singleton)
```

### Type Mismatches

Calling `Get` or `GetOrDo` for a result type that doesn't match the previously stored result type for a `Singleton` will result in an `ErrTypeMismatch`:

```go
singlet.GetOrDo(s, cache.New[int])
singlet.Get[string](s) // Returns ErrTypeMismatch
```

## License

Copyright Jonathan Halterman. Released under the [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html).