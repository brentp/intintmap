// Package intintmap is a fast int64 key -> int64 value map.
//
// It is copied nearly verbatim from http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/
package intintmap

import "math"

const INT_PHI = 0x9E3779B9
const FREE_KEY = 0
const NO_VALUE = 0

func phiMix(x int64) int64 {
	h := x * INT_PHI
	return h ^ (h >> 16)
}

// Map is a map-like data-structure for int64s
type Map struct {
	data       []int64
	fillFactor float64
	threshold  int

	mask  int64
	mask2 int64

	hasFreeKey bool
	freeVal    int64
	size       int
}

func nextPowerOf2(x uint32) uint32 {
	if x == 0 {
		return 1
	}
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return (x | x>>32) + 1
}

func arraySize(exp int, fill float64) int {
	s := nextPowerOf2(uint32(math.Ceil(float64(exp) / fill)))
	if s < 2 {
		s = 2
	}
	return int(s)
}

// New returns a map initialized with n spaces and uses the stated fillFactor.
// The map will grow as needed.
func New(n int, fillFactor float64) *Map {
	capacity := arraySize(n, fillFactor)
	return &Map{
		data:       make([]int64, 2*capacity),
		fillFactor: fillFactor,
		threshold:  int(math.Floor(float64(capacity) * fillFactor)),
		mask:       int64(capacity - 1),
		mask2:      int64(2*capacity - 1),
	}
}

// Get returns the value or NO_VALUE if the key is not found.
func (m *Map) Get(key int64) int64 {

	if key == FREE_KEY {
		if m.hasFreeKey {
			return m.freeVal
		}
		return NO_VALUE
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == key {
		return m.data[ptr+1]
	}

	for k != key && k != FREE_KEY {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]
	}

	if k == FREE_KEY {
		return NO_VALUE
	}
	return m.data[ptr+1]
}

// Put adds val to the map under the specified key and returns the old value in that key.
func (m *Map) Put(key int64, val int64) int64 {
	if key == FREE_KEY {
		ret := m.freeVal
		if !m.hasFreeKey {
			m.size += 1
		}
		m.hasFreeKey = true
		m.freeVal = val
		return ret
	}

	ptr := (phiMix(key) & m.mask) << 1
	k := m.data[ptr]

	if k == FREE_KEY { // end of chain
		m.data[ptr+1] = val
		m.data[ptr] = key
		if m.size >= m.threshold {
			m.rehash()
		} else {
			m.size += 1
		}
		return NO_VALUE
	} else if k == key {
		ret := m.data[ptr+1]
		m.data[ptr+1] = val
		return ret
	}

	for {
		ptr = (ptr + 2) & m.mask2
		k = m.data[ptr]

		if k == FREE_KEY {
			m.data[ptr+1] = val
			m.data[ptr] = key
			if m.size >= m.threshold {
				m.rehash()
			} else {
				m.size += 1
			}
			return NO_VALUE
		} else if k == key {
			ret := m.data[ptr+1]
			m.data[ptr+1] = val
			return ret
		}
	}

}

func (m *Map) rehash() {
	newCapacity := len(m.data) * 2
	m.threshold = int(math.Floor(float64(newCapacity/2) * m.fillFactor))
	m.mask = int64(newCapacity/2 - 1)
	m.mask2 = int64(newCapacity - 1)

	data := make([]int64, len(m.data))
	copy(data, m.data)

	m.data = make([]int64, newCapacity)

	for i := 0; i < len(data); i += 2 {
		o := data[i]
		if o != FREE_KEY {
			m.Put(o, data[i+1])
		}
	}
}
