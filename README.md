Fast int64 -> int64 hash in golang.

# intintmap
--
    import "github.com/brentp/intintmap"

Package intintmap is a fast int64 key -> int64 value map.

It is copied nearly verbatim from
http://java-performance.info/implementing-world-fastest-java-int-to-int-hash-map/

It is 2-4X faster than the builtin map:
```
BenchmarkIntIntMapFill    2000   1080930 ns/op
BenchmarkStdMapFill       1000   1928913 ns/op
BenchmarkIntIntMapGet     20     96848089 ns/op
BenchmarkStdMapGet        5      246061427 ns/op
```


## Usage


```go
m := intintmap.New(32768, 0.6)
m.Put(int64(1234), int64(-222))
m.Put(int64(123), int64(33))

m.Get(int64(222))
m.Get(int64(333))
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
