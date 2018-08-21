Fast int64 -> int64 hash in golang.

[![GoDoc](https://godoc.org/github.com/brentp/intintmap?status.svg)](https://godoc.org/github.com/brentp/intintmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/brentp/intintmap)](https://goreportcard.com/report/github.com/brentp/intintmap)

# intintmap

    import "github.com/brentp/intintmap"

Package intintmap is a fast int64 key -> int64 value map.

It is copied nearly verbatim from
http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/ .

It interleaves keys and values in the same underlying array to improve locality.

It is 2-5X faster than the builtin map:
```
BenchmarkIntIntMapFill                 	      10	 158436598 ns/op
BenchmarkStdMapFill                    	       5	 312135474 ns/op
BenchmarkIntIntMapGet10PercentHitRate  	    5000	    243108 ns/op
BenchmarkStdMapGet10PercentHitRate     	    5000	    268927 ns/op
BenchmarkIntIntMapGet100PercentHitRate 	     500	   2249349 ns/op
BenchmarkStdMapGet100PercentHitRate    	     100	  10258929 ns/op
```

**note** it currently returns 0 for missing keys. Should probably make Get() return (int, error) so we can
differentiate (or use math.MinUint32 if that hurts performance).

## Usage

```go
m := intintmap.New(32768, 0.6)
m.Put(int64(1234), int64(-222))
m.Put(int64(123), int64(33))

m.Get(int64(222))
m.Get(int64(333))

m.Del(int64(222))
m.Del(int64(333))

m.Size()

for k := range m.Keys() {
    fmt.Printf("key: %d\n", k)
}

for kv := range m.Items() {
    fmt.Printf("key: %d, value: %d\n", kv[0], kv[1])
}
```

#### type Map

```go
type Map struct {
}
```

Map is a map-like data-structure for int64s

#### func  New

```go
func New(n int, fillFactor float64) *Map
```
New returns a map initialized with n spaces and uses the stated fillFactor. The
map will grow as needed.

#### func (*Map) Get

```go
func (m *Map) Get(key int64) int64
```
Get returns the value or NO_VALUE if the key is not found.

#### func (*Map) Put

```go
func (m *Map) Put(key int64, val int64) int64
```
Put adds val to the map under the specified key and returns the old value in
that key.

#### func (*Map) Del

```go
func (m *Map) Del(key int64) int64
```
Del deletes a key and its value, and returns the value or NO_VALUE if the
key is not found.

#### func (*Map) Keys

```go
func (m *Map) Keys() chan int64
```
Keys returns a channel for iterating all keys.

#### func (*Map) Items

```go
func (m *Map) Items() chan [2]int64
```
Items returns a channel for iterating all key-value pairs.


#### func (*Map) Size

```go
func (m *Map) Size() int
```
Size returns size of the map.
