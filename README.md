# xmap

![Build Status](https://github.com/icza/xmap/actions/workflows/go.yml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/icza/xmap.svg)](https://pkg.go.dev/github.com/icza/xmap)

Package `xmap` provides a hashmap implementation that uses custom, user-defined
`eq()` and `hash()` functions for equality checks and hash calculation.
The implementation is not safe for concurrent use.

Performance depends on how fast the provided `eq()` and `hash()` functions are.
The implementation minimizes calls on the `hash()` function by caching hashes.
Cached hashes are also used to minimize `eq()` calls.
As a result, for simple types `xmap` outperforms Go's builtin map significantly
(see `BenchmarkIntMap`). For complex types and inefficient `hash()` functions, performance
will be worse.

This package is experimental and is not production ready yet.

## Example

	type Key struct {
		a *int
	}

	func NewKey(a int) Key           { return Key{a: &a} }
	func myEquals(x Key, y Key) bool { return *x.a == *y.a }
	func myHash(k Key) uint32        { return uint32(*k.a) }

	func Example() {
		m := xmap.New[Key, int](myEquals, myHash)
		// Add some entries:
		for i := 0; i < 5; i++ {
			m.Set(NewKey(i), i+10)
		}
		// Get and remove:
		if v, ok := m.GetOK(NewKey(2)); ok {
			m.Delete(NewKey(v - 10 + 1)) // This will remove key=3
		}

		fmt.Println("Range:")
		m.Range(func(k Key, v int) bool {
			fmt.Println("\tentry:", *k.a, "=", v)
			return true
		})

		// Output:
		// Range:
		// 	entry: 0 = 10
		// 	entry: 1 = 11
		// 	entry: 2 = 12
		// 	entry: 4 = 14
	}
