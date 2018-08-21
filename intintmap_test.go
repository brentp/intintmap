package intintmap

import (
	"testing"
)

func TestMapSimple(t *testing.T) {
	m := New(10, 0.99)
	var i int64
	var v int64
	var ok bool

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 2; i < 20000; i += 2 {
		m.Put(i, i)
	}
	for i = 2; i < 20000; i += 2 {
		if v, ok = m.Get(i); !ok || v != i {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	if m.Size() != int(20000/2)-1 {
		t.Errorf("size (%d) is not right, should be %d", m.Size(), int(20000/2)-1)
	}

	// --------------------------------------------------------------------
	// Keys()

	m0 := make(map[int64]int64, 1000)
	for i = 2; i < 20000; i += 2 {
		m0[i] = i
	}
	n := len(m0)

	for k := range m.Keys() {
		m0[k] = -k
	}
	if n != len(m0) {
		t.Errorf("get unexpected more keys")
	}

	for k, v := range m0 {
		if k != -v {
			t.Errorf("didn't get expected changed value")
		}
	}

	// --------------------------------------------------------------------
	// Items()

	m0 = make(map[int64]int64, 1000)
	for i = 2; i < 20000; i += 2 {
		m0[i] = i
	}
	n = len(m0)

	for kv := range m.Items() {
		m0[kv[0]] = -kv[1]
		if kv[0] != kv[1] {
			t.Errorf("didn't get expected key-value pair")
		}
	}
	if n != len(m0) {
		t.Errorf("get unexpected more keys")
	}

	for k, v := range m0 {
		if k != -v {
			t.Errorf("didn't get expected changed value")
		}
	}

	// --------------------------------------------------------------------
	// Del()

	for i = 2; i < 20000; i += 2 {
		m.Del(i)
	}
	for i = 2; i < 20000; i += 2 {
		if _, ok = m.Get(i); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

	// --------------------------------------------------------------------
	// Put() and Get()

	for i = 2; i < 20000; i += 2 {
		m.Put(i, i*2)
	}
	for i = 2; i < 20000; i += 2 {
		if v, ok = m.Get(i); !ok || v != i*2 {
			t.Errorf("didn't get expected value")
		}
		if _, ok = m.Get(i + 1); ok {
			t.Errorf("didn't get expected 'not found' flag")
		}
	}

}

func TestMap(t *testing.T) {
	m := New(10, 0.6)
	var ok bool
	var v int64

	step := int64(61)

	var i int64
	m.Put(0, 12345)
	for i = 1; i < 100000000; i += step {
		m.Put(i, i+7)
		m.Put(-i, i-7)

		if v, ok = m.Get(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Get(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}
	}
	for i = 1; i < 100000000; i += step {
		if v, ok = m.Get(i); !ok || v != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, v)
		}
		if v, ok = m.Get(-i); !ok || v != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, v)
		}

		for j := i + 1; j < i+step; j++ {
			if v, ok = m.Get(j); ok {
				t.Errorf("expected 'not found' flag for %d, found %d", j, v)
			}
		}
	}

	if v, ok = m.Get(0); !ok || v != 12345 {
		t.Errorf("expected 12345 for key 0")
	}
}

const MAX = 999999999
const STEP = 9534

func fillIntIntMap(m *Map) {
	var j int64
	for j = 0; j < MAX; j += STEP {
		m.Put(j, -j)
		for k := j; k < j+16; k++ {
			m.Put(k, -k)
		}

	}
}

func fillStdMap(m map[int64]int64) {
	var j int64
	for j = 0; j < MAX; j += STEP {
		m[j] = -j
		for k := j; k < j+16; k++ {
			m[k] = -k
		}
	}
}

func BenchmarkIntIntMapFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := New(2048, 0.60)
		fillIntIntMap(m)
	}
}

func BenchmarkStdMapFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int64]int64, 2048)
		fillStdMap(m)
	}
}

func BenchmarkIntIntMapGet10PercentHitRate(b *testing.B) {
	var j, k, v, sum int64
	var ok bool
	m := New(2048, 0.60)
	fillIntIntMap(m)
	for i := 0; i < b.N; i++ {
		sum = int64(0)
		for j = 0; j < MAX; j += STEP {
			for k = j; k < 10; k++ {
				if v, ok = m.Get(k); ok {
					sum += v
				}
			}
		}
		//log.Println("int int sum:", sum)
	}
}

func BenchmarkStdMapGet10PercentHitRate(b *testing.B) {
	var j, k, v, sum int64
	var ok bool
	m := make(map[int64]int64, 2048)
	fillStdMap(m)
	for i := 0; i < b.N; i++ {
		sum = int64(0)
		for j = 0; j < MAX; j += STEP {
			for k = j; k < 10; k++ {
				if v, ok = m[k]; ok {
					sum += v
				}
			}
		}
		//log.Println("map sum:", sum)
	}
}

func BenchmarkIntIntMapGet100PercentHitRate(b *testing.B) {
	var j, v, sum int64
	var ok bool
	m := New(2048, 0.60)
	fillIntIntMap(m)
	for i := 0; i < b.N; i++ {
		sum = int64(0)
		for j = 0; j < MAX; j += STEP {
			if v, ok = m.Get(j); ok {
				sum += v
			}
		}
		//log.Println("int int sum:", sum)
	}
}

func BenchmarkStdMapGet100PercentHitRate(b *testing.B) {
	var j, v, sum int64
	var ok bool
	m := make(map[int64]int64, 2048)
	fillStdMap(m)
	for i := 0; i < b.N; i++ {
		sum = int64(0)
		for j = 0; j < MAX; j += STEP {
			if v, ok = m[j]; ok {
				sum += v
			}
		}
		//log.Println("map sum:", sum)
	}
}
