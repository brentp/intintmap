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
const STEP = 95346

func fillIntIntMap(m *Map) {
	var j int64
	for j = 0; j < MAX; j += STEP {
		m.Put(j, -j)
	}

}

func fillStdMap(m map[int64]int64) {
	var j int64
	for j = 0; j < MAX; j += STEP {
		m[j] = -j
	}
}

func BenchmarkIntIntMapFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := New(2048, 0.60)
		fillIntIntMap(m)
		fillIntIntMap(m)
	}
}

func BenchmarkStdMapFill(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int64]int64, 2048)
		fillStdMap(m)
		fillStdMap(m)
	}
}

func BenchmarkIntIntMapGet(b *testing.B) {
	var j int64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sum := int64(0)
		m := New(2048, 0.60)
		fillIntIntMap(m)
		b.StartTimer()
		for j = 0; j < MAX; j += 100 {
			if v := m.Get(j); v != NO_VALUE {
				sum += v
			}
		}
		//log.Println("int int sum:", sum)
	}
}

func BenchmarkStdMapGet(b *testing.B) {
	var j int64
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		sum := int64(0)
		m := make(map[int64]int64, 2048)
		fillStdMap(m)
		b.StartTimer()
		for j = 0; j < MAX; j += 100 {
			if v, ok := m[j]; ok {
				sum += v
			}
		}
		//log.Println("map sum:", sum)
	}
}
