package intintmap

import "testing"

func TestMap(t *testing.T) {
	m := New(10, 0.6)

	step := int64(61)

	var i int64
	m.Put(0, 12345)
	for i = 1; i < 100000000; i += step {
		m.Put(i, i+7)
		m.Put(-i, i-7)

		if m.Get(i) != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, m.Get(i))
		}
		if m.Get(-i) != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, m.Get(-i))
		}
	}
	for i = 1; i < 100000000; i += step {
		if m.Get(i) != i+7 {
			t.Errorf("expected %d as value for key %d, got %d", i+7, i, m.Get(i))
		}
		if m.Get(-i) != i-7 {
			t.Errorf("expected %d as value for key %d, got %d", i-7, -i, m.Get(-i))
		}

		for j := i + 1; j < i+step; j++ {
			if m.Get(j) != NO_VALUE {
				t.Errorf("expected empty value for %d, found %d", j, m.Get(j))
			}
		}
	}

	if m.Get(0) != 12345 {
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
	var j int64
	m := New(2048, 0.60)
	fillIntIntMap(m)
	for i := 0; i < b.N; i++ {
		sum := int64(0)
		for j = 0; j < MAX; j += STEP {
			for k := j; k < 10; k++ {
				if v := m.Get(k); v != NO_VALUE {
					sum += v
				}
			}
		}
		//log.Println("int int sum:", sum)
	}
}

func BenchmarkStdMapGet10PercentHitRate(b *testing.B) {
	var j int64
	m := make(map[int64]int64, 2048)
	fillStdMap(m)
	for i := 0; i < b.N; i++ {
		sum := int64(0)
		for j = 0; j < MAX; j += STEP {
			for k := j; k < 10; k++ {
				if v, ok := m[k]; ok {
					sum += v
				}
			}
		}
		//log.Println("map sum:", sum)
	}
}

func BenchmarkIntIntMapGet100PercentHitRate(b *testing.B) {
	var j int64
	m := New(2048, 0.60)
	fillIntIntMap(m)
	for i := 0; i < b.N; i++ {
		sum := int64(0)
		for j = 0; j < MAX; j += STEP {
			if v := m.Get(j); v != NO_VALUE {
				sum += v
			}
		}
		//log.Println("int int sum:", sum)
	}
}

func BenchmarkStdMapGet100PercentHitRate(b *testing.B) {
	var j int64
	m := make(map[int64]int64, 2048)
	fillStdMap(m)
	for i := 0; i < b.N; i++ {
		sum := int64(0)
		for j = 0; j < MAX; j += STEP {
			if v, ok := m[j]; ok {
				sum += v
			}
		}
		//log.Println("map sum:", sum)
	}
}
