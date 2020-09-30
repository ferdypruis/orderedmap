# orderedmap
An ordered Go `map[string]string`, to for example generate or parse a JSON object with sorted keys.

Despite objects being an unordered set of name/value pairs in JSON, there are some implementations
that do require a particular key order.
 
## Examples
Marshal an ordered set
```go
var m StringMap
m.Set("first", "1")
m.Set("second", "2")
m.Set("third", "3")

out, _ := json.Marshal(m)
fmt.Println(string(out))

// Output:
// {"first":"1","second":"2","third":"3"}
```

Unmarshal an object and maintain key order
```go
data := []byte(`{"first":"1","second":"2","third":"3"}`)

var m StringMap
_ = json.Unmarshal(data, &m)

for _, k := range m.Keys() {
    v, _ := m.Value(k)
    fmt.Println(k, "=", v)
}

// Output:
// first = 1
// second = 2
// third = 3
```

`StringMap` implements `sort.Interface`
```go
var m StringMap
m.Set("first", "1")
m.Set("second", "2")
m.Set("third", "3")

// Reverse sort values
sort.Sort(sort.Reverse(m))

for _, k := range m.Keys() {
    v, _ := m.Value(k)
    fmt.Println(k, "=", v)
}

// Output:
// third = 3
// second = 2
// first = 1
```