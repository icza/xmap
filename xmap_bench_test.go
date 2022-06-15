package xmap

import (
	"encoding/binary"
	"fmt"
	"hash/maphash"
	"testing"
)

var int1, int2 int

func BenchmarkIntMap(b *testing.B) {
	for entries := 10; entries <= 1_000_000; entries *= 10 {
		b.Run(fmt.Sprintf("Go int map n=%d", entries), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := map[int]int{}
				for j := 0; j < entries; j++ {
					m[j] = j
					int1 = m[j]
				}
				for k, v := range m {
					int1, int2 = k, v
				}
				for j := 0; j < entries; j++ {
					int1 = m[j]
					delete(m, j)
					int1 = m[j]
				}
			}
		})

		b.Run(fmt.Sprintf("XMAP int   n=%d", entries), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := New[int, int](intEq, intHash)
				for j := 0; j < entries; j++ {
					m.Set(j, j)
					int1 = m.Get(j)
				}
				m.Range(func(k, v int) bool {
					int1, int2 = k, v
					return true
				})
				for j := 0; j < entries; j++ {
					int1 = m.Get(j)
					m.Delete(j)
					int1 = m.Get(j)
				}
			}
		})
	}
}

type person struct {
	name string
	age  int
	num  float64
}

var personSeed = maphash.MakeSeed()

func personHash(p person) uint32 {
	var h = new(maphash.Hash)
	h.SetSeed(personSeed)
	h.WriteString(p.name)
	if err := binary.Write(h, binary.LittleEndian, int64(p.age)); err != nil {
		panic(err)
	}
	if err := binary.Write(h, binary.LittleEndian, p.num); err != nil {
		panic(err)
	}

	return uint32(h.Sum64())
}

func personEq(p1, p2 person) bool {
	return p1 == p2
}

var person1, person2 person

var names []string

func init() {
	for i := 0; i < 100; i++ {
		names = append(names, fmt.Sprintf("Bob %d", i))
		names = append(names, fmt.Sprintf("%d Alice", i*2))
	}
}

func getPerson(j int) person {
	return person{
		name: names[j%len(names)],
		age:  j,
		num:  float64(j),
	}
}

func BenchmarkPersonMap(b *testing.B) {
	for entries := 10; entries <= 1_000_000; entries *= 10 {
		b.Run(fmt.Sprintf("Go person map n=%d", entries), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := map[person]person{}
				for j := 0; j < entries; j++ {
					p := getPerson(j)
					m[p] = p
					person1 = m[p]
				}
				for k, v := range m {
					person1, person2 = k, v
				}
				for j := 0; j < entries; j++ {
					p := getPerson(j)
					person1 = m[p]
					delete(m, p)
					person1 = m[p]
				}
			}
		})

		b.Run(fmt.Sprintf("XMAP person   n=%d", entries), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				m := New[person, person](personEq, personHash)
				for j := 0; j < entries; j++ {
					p := getPerson(j)
					m.Set(p, p)
					person1 = m.Get(p)
				}
				m.Range(func(k, v person) bool {
					person1, person2 = k, v
					return true
				})
				for j := 0; j < entries; j++ {
					p := getPerson(j)
					person1 = m.Get(p)
					m.Delete(p)
					person1 = m.Get(p)
				}
			}
		})
	}
}
