package orderedmap_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	. "github.com/ferdypruis/orderedmap"
)

func TestStringMap(t *testing.T) {
	data := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	// This key should be overwritten
	stringmap.Set("key one", "value ?")
	for _, d := range data {
		stringmap.Set(d.k, d.v)
	}

	// Assert contents of stringmap
	keys := stringmap.Keys()
	if len(keys) != len(data) {
		t.Errorf("expected %d keys, got %d; %#v", len(data), stringmap.Len(), keys)
	}

	for i, key := range keys {
		if key != data[i].k {
			t.Errorf("expected key %d to be %q, got %q", i, data[i].k, key)
		} else if value, ok := stringmap.Value(key); !ok {
			t.Errorf("expected value for key %q to exist", key)
		} else if value != data[i].v {
			t.Errorf("expected value for key %q to be %q, got %q", key, data[i].v, key)
		}
	}

	if value, ok := stringmap.Value("notexist"); ok {
		t.Errorf("expected value for key %q not to exist, got %q", "notexist", value)
	}
}

func TestStringmap_MarshalJSON(t *testing.T) {
	var stringmap StringMap
	stringmap.Set("key one", "value 1")
	stringmap.Set("otherkey", "val2")
	stringmap.Set("key3", "a third value")

	actually, err := json.Marshal(stringmap)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`{"key one":"value 1","otherkey":"val2","key3":"a third value"}`)
	if !bytes.Equal(actually, expected) {
		t.Errorf("expected json %s, got %s", expected, actually)
	}
}

func TestStringmap_MarshalJSONEmpty(t *testing.T) {
	var stringmap StringMap
	actually, err := json.Marshal(stringmap)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(`{}`)
	if !bytes.Equal(actually, expected) {
		t.Errorf("expected json %s, got %s", expected, actually)
	}
}

func TestStringmap_UnmarshalJSON(t *testing.T) {
	expected := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	err := json.Unmarshal([]byte(`{"key one":"value 1","otherkey":"val2","key2":"a third value"}`), &stringmap)
	if err != nil {
		t.Fatal(err)
	}

	if stringmap.Len() != len(expected) {
		t.Errorf("expected %d items, got %d", len(expected), stringmap.Len())
	}
	for i, key := range stringmap.Keys() {
		if key != expected[i].k {
			t.Errorf("expected item %d to have key %q, got %q", i, expected[i].k, key)
		}
		if value, _ := stringmap.Value(key); value != expected[i].v {
			t.Errorf("expected item %d to have value %q, got %q", i, expected[i].v, value)
		}
	}
}

func TestStringmap_UnmarshalJSONErrors(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
	}{
		{"empty input", []byte("")},
		{"json null value", []byte("null")},
		{"json string value", []byte(`"hello"`)},
		{"invalid key type", []byte(`{231:"no"}`)},
		{"error value", []byte(`{"nietes":welles}`)},
		{"invalid value type", []byte(`{"number":231}`)},
		{"invalid end of object", []byte(`{"key": "val" `)},
		{"trailing data", []byte(`{"key": "val" },`)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stringmap StringMap
			if err := stringmap.UnmarshalJSON(test.input); err == nil {
				t.Errorf("expected error")
			}
		})
	}
}

func TestStringmapSort(t *testing.T) {
	data := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	for _, d := range data {
		stringmap.Set(d.k, d.v)
	}

	// Regular sort of values
	sort.Sort(stringmap)

	// sort order is digit, lowercase, uppercase
	expected := []struct {
		k string
		v string
	}{
		{"key2", "a third value"},
		{"otherkey", "val2"},
		{"key one", "value 1"},
	}

	for i, key := range stringmap.Keys() {
		if key != expected[i].k {
			t.Errorf("expected item %d to have key %q, got %q", i, expected[i].k, key)
		}
		if value, _ := stringmap.Value(key); value != expected[i].v {
			t.Errorf("expected item %d to have value %q, got %q", i, expected[i].v, value)
		}
	}
}

func TestStringmap_Sort(t *testing.T) {
	data := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	for _, d := range data {
		stringmap.Set(d.k, d.v)
	}

	// Sort by the length of the value
	stringmap.Sort(func(s, t string) bool {
		return len(s) < len(t)
	})

	// sort order is digit, lowercase, uppercase
	expected := []struct {
		k string
		v string
	}{
		{"otherkey", "val2"},
		{"key one", "value 1"},
		{"key2", "a third value"},
	}

	for i, key := range stringmap.Keys() {
		if key != expected[i].k {
			t.Errorf("expected item %d to have key %q, got %q", i, expected[i].k, key)
		}
		if value, _ := stringmap.Value(key); value != expected[i].v {
			t.Errorf("expected item %d to have value %q, got %q", i, expected[i].v, value)
		}
	}
}

func TestStringmap_SortKeys(t *testing.T) {
	data := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	for _, d := range data {
		stringmap.Set(d.k, d.v)
	}

	// Sort by the length of the key
	stringmap.SortKeys(func(s, t string) bool {
		return len(s) < len(t)
	})

	// sort order is digit, lowercase, uppercase
	expected := []struct {
		k string
		v string
	}{
		{"key2", "a third value"},
		{"key one", "value 1"},
		{"otherkey", "val2"},
	}

	for i, key := range stringmap.Keys() {
		if key != expected[i].k {
			t.Errorf("expected item %d to have key %q, got %q", i, expected[i].k, key)
		}
		if value, _ := stringmap.Value(key); value != expected[i].v {
			t.Errorf("expected item %d to have value %q, got %q", i, expected[i].v, value)
		}
	}
}

// TestStringMap_KeysImmutable asserts we can not manipulate the keys
func TestStringMap_KeysImmutable(t *testing.T) {
	data := []struct {
		k string
		v string
	}{
		{"key one", "value 1"},
		{"otherkey", "val2"},
		{"key2", "a third value"},
	}

	var stringmap StringMap
	for _, d := range data {
		stringmap.Set(d.k, d.v)
	}

	keys := stringmap.Keys()
	keys[0] = "fu"
	keys[1] = "bar"

	// Now check and see stringmap has not changed
	for i, key := range stringmap.Keys() {
		if key != data[i].k {
			t.Errorf("expected key %d to be %q, got %q", i, data[i].k, key)
		}
	}
}

func ExampleStringMap_MarshalJSON() {
	var m StringMap
	m.Set("first", "1")
	m.Set("second", "2")
	m.Set("third", "3")

	out, _ := json.Marshal(m)
	fmt.Println(string(out))

	// Output:
	// {"first":"1","second":"2","third":"3"}
}

func ExampleStringMap_UnmarshalJSON() {
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
}

func ExampleStringMap_Sortable() {
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
}
