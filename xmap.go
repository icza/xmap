/*

Package xmap provides a hashmap implementation that uses custom, user-provided
eq() and hash() functions for equality checks and hash calculation.
The implementation is not safe for concurrent use.

Performance depends on how fast the provided eq() and hash() functions are.
The implementation minimizes calls on the hash() function by caching hashes.
As a result, for simple types xmap outperforms Go's builtin map significantly
(see BenchmarkIntMap). For complex types and inefficient hash() functions, performance
will be worse.

This package is experimental and is not yet production ready.

*/
package xmap

import (
	"math"
)

// Map is a hashmap implementation that uses custom eq() and hash() functions.
type Map[K, V any] struct {
	eq   func(key1, key2 K) bool
	hash func(key K) uint32
	cfg  Config

	maxSize int // maxSize to hold without rehashing
	minSize int // minSize to hold without rehashing

	size    int // current size
	buckets []*entry[K, V]
}

// New creates a new Map, using the given eq and hash functions
// for equality checks and hash calculations.
func New[K, V any](
	eq func(key1, key2 K) bool,
	hash func(key K) uint32,
) *Map[K, V] {
	return NewConfig[K, V](eq, hash, defaultConfig)
}

func NewConfig[K, V any](
	eq func(key1, key2 K) bool,
	hash func(key K) uint32,
	cfg *Config,
) *Map[K, V] {

	m := &Map[K, V]{
		eq:   eq,
		hash: hash,
		cfg:  *cfg.setDefaults(),
	}

	return m
}

type entry[K, V any] struct {
	key   K
	value V
	hash  uint32 // Cached hash of key
	next  *entry[K, V]
}

// Get returns the value associated with key.
// Returns the zero value of the value type if not found.
func (m *Map[K, V]) Get(key K) (value V) {
	value, _ = m.GetOK(key)
	return
}

// GetOK returns the value associated with key.
// If the key is found, ok will be true, false otherwise.
// Returns the zero value of the value type if not found.
func (m *Map[K, V]) GetOK(key K) (value V, ok bool) {
	if m.size == 0 {
		return
	}

	hash := m.hash(key)
	idx := m.hashToBucketIdx(hash)

	for e := m.buckets[idx]; e != nil; e = e.next {
		if hash == e.hash && // Quick check
			m.eq(key, e.key) {
			return e.value, true
		}
	}

	return
}

// Set sets the given value for the given key.
// If the key is already in the map, updates the value.
func (m *Map[K, V]) Set(key K, value V) {
	if m.size == 0 {
		m.SetCap(m.cfg.InitialCap)
	}

	hash := m.hash(key)
	idx := m.hashToBucketIdx(hash)

	if e := m.buckets[idx]; e == nil {
		// This is the first in bucket
		m.buckets[idx] = &entry[K, V]{
			key:   key,
			value: value,
			hash:  hash,
		}
	} else {
		// Find in bucket, and update if found, else append
		for {
			if hash == e.hash && // Quick check
				m.eq(key, e.key) {
				// Found, just update:
				e.value = value
				return
			}
			if e.next == nil {
				// Not found, append
				e.next = &entry[K, V]{
					key:   key,
					value: value,
					hash:  hash,
				}
				break
			}
			e = e.next
		}
	}

	// A new entry was just inserted:
	m.changeSize(1)
}

// Delete removes the key and its value from the map.
func (m *Map[K, V]) Delete(key K) {
	if m.size == 0 {
		return
	}

	hash := m.hash(key)
	idx := m.hashToBucketIdx(hash)

	var prev *entry[K, V]

	for e := m.buckets[idx]; e != nil; e = e.next {
		if hash == e.hash && // Quick check
			m.eq(key, e.key) {

			// Found, remove
			if prev == nil {
				m.buckets[idx] = e.next // The next is the new first
			} else {
				prev.next = e.next // "Link through" the deleted entry
			}

			// An entry was just removed:
			m.changeSize(-1)
			return
		}
		prev = e
	}
}

// Range ranges over the entries of the map, calling f for each key-value pairs.
// If f returns false, Range will return without further calls to f.
//
// If new entries are added during the run of Range, they may or may not be visited.
// If entries are removed during the run of Range, it might result in other unvisited entries to be skipped.
// If the map is structurally modified during the run of Range (e.g. a rehash happens
// due to entry addition or removal or an explicit call to SetCap()), further entry visits will be undefined
// (it may be already visited entries will be visited again or unvisited entries will be skipped).
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	for _, e := range m.buckets {
		for ; e != nil; e = e.next {
			if !f(e.key, e.value) {
				return
			}
		}
	}
}

func (m *Map[K, V]) hashToBucketIdx(hash uint32) uint32 {
	return hash % uint32(len(m.buckets))
}

// changeSize changes the current size and checks if its in the allowed [minSize..maxSize] range.
// Performs rehashing to grow / shrink as needed.
func (m *Map[K, V]) changeSize(delta int) {
	m.size += delta

	if delta > 0 { // Size increased
		if m.size > m.maxSize {
			newSize := int(float64(m.maxSize) * m.cfg.ChangeFactor)
			m.SetCap(newSize)
		}
	} else { // Size decreased
		if m.size < m.minSize && m.maxSize > m.cfg.InitialCap { // Don't go below the initial capacity
			newSize := int(float64(m.maxSize) / m.cfg.ChangeFactor)
			if newSize < m.cfg.InitialCap { // Don't go below the initial capacity
				newSize = m.cfg.InitialCap
			}
			m.SetCap(newSize)
		}
	}
}

// Len returns the map's length (the number of entries it holds).
func (m *Map[K, V]) Len() int {
	return m.size
}

// Cap returns the map's current capacity.
// The capacity is the number of entries the map can hold without causing a rehashing.
func (m *Map[K, V]) Cap() int {
	return m.maxSize
}

// SetCap sets the map's capacity.
// The capacity is the number of entries the map can hold without causing a rehashing.
// If the capacity is bigger than the current capacity, the map will grow. If it's smaller, the map will shrink.
//
// If the requested capacity would be insufficient to hold current entries, the minimum required capacity
// will be used instead (which is the map's current size).
func (m *Map[K, V]) SetCap(capacity int) {
	if capacity < m.size {
		capacity = m.size
	}

	if capacity == m.maxSize {
		return // We already have the requested capacity.
	}

	m.maxSize = capacity
	numBuckets := int(math.Ceil(float64(m.maxSize) / m.cfg.GrowLoadLimit))
	m.minSize = int(math.Ceil(float64(numBuckets) * m.cfg.ShrinkLoadLimit))

	m.rehash(numBuckets)
}

// rehash performs a rehashing.
func (m *Map[K, V]) rehash(numBuckets int) {
	buckets := m.buckets
	m.buckets = make([]*entry[K, V], numBuckets)

	// Build new buckets, keep existing entries but clear next fields.
	for _, e := range buckets {
		for e != nil {
			next := e.next
			e.next = nil

			idx := m.hashToBucketIdx(e.hash)
			if last := m.buckets[idx]; last == nil {
				m.buckets[idx] = e // This is the first in bucket
			} else {
				for last.next != nil {
					last = last.next
				}
				last.next = e // Append to last
			}

			e = next
		}
	}
}

// Clone returns a clone of this map, holding the same entries, having the same config and capacity.
func (m *Map[K, V]) Clone() *Map[K, V] {
	m2 := new(Map[K, V])
	*m2 = *m

	// Deepcopy buckets
	m2.buckets = make([]*entry[K, V], len(m.buckets))
	for i, e := range m.buckets {
		var prev *entry[K, V]

		for ; e != nil; e = e.next {
			e2 := new(entry[K, V])
			*e2 = *e
			e2.next = nil

			if prev == nil {
				m2.buckets[i] = e2 // First in bucket
			} else {
				prev.next = e2 // Link to this entry
			}

			prev = e2
		}
	}

	return m2
}
