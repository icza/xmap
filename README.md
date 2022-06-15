# xmap

![Build Status](https://github.com/icza/xmap/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/icza/xmap.svg)](https://pkg.go.dev/github.com/icza/xmap)

Package `xmap` provides a hashmap implementation that uses custom, user-provided
`eq()` and `hash()` functions for equality checks and hash calculation.
The implementation is not safe for concurrent use.

Performance depends on how fast the provided `eq()` and `hash()` functions are.
The implementation minimizes calls on the `hash()` function by caching hashes.
Cached hashes are also used to minimize `eq()` calls.
As a result, for simple types `xmap` outperforms Go's builtin map significantly
(see `BenchmarkIntMap`). For complex types and inefficient `hash()` functions, performance
will be worse.

This package is experimental and is not production ready yet.

