package xmap

import (
	"testing"

	"github.com/icza/mighty"
)

func intHash(i int) uint32 {
	return uint32(i)
}

func intEq(i, j int) bool {
	return i == j
}

func TestNormal(t *testing.T) {
	eq := mighty.Eq(t)

	m := New[int, int](intEq, intHash)

	eq(0, m.Len())
	_, ok := m.GetOK(1)
	eq(false, ok)

	m.Set(1, 11)
	eq(1, m.Len())
	i, ok := m.GetOK(1)
	eq(true, ok)
	eq(11, i)
}

func TestCapacity(t *testing.T) {
	eq := mighty.Eq(t)

	cfg := &Config{
		InitialCap:      10,
		ChangeFactor:    2,
		GrowLoadLimit:   0.8,
		ShrinkLoadLimit: 0.25,
	}
	m := NewConfig[int, int](intEq, intHash, cfg)

	eq(0, m.Len())
	eq(0, m.Cap())

	for i := 1; i < 30; i++ {
		m.Set(i, i)
		eq(i, m.Len())

		switch {
		case i >= 21: // At len=21 capacity increases to 40
			eq(40, m.Cap())
		case i >= 11: // At len=11 capacity increases to 20
			eq(20, m.Cap())
		case i >= 1: // At len=1 capacity increases to 10
			eq(10, m.Cap())
		}
	}

	for i := 30 - 1; i > 0; i-- {
		m.Delete(i)
		eq(i-1, m.Len())

		switch {
		case i > 13: // At len=13 capacity decreases to 20
			eq(40, m.Cap())
		case i > 7: // At len=7 capacity decreases to 10
			eq(20, m.Cap())
		default: // Capacity doesn't fall below the initial 10
			eq(10, m.Cap())
		}
	}
}
