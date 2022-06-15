package xmap_test

import (
	"fmt"

	"github.com/icza/xmap"
)

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
