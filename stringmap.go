package orderedmap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
)

var _ json.Marshaler = (*StringMap)(nil)
var _ json.Unmarshaler = (*StringMap)(nil)
var _ sort.Interface = (*StringMap)(nil)

// StringMap represents a map of string key/value pairs which maintains its order when marshaled to/from JSON
// Like the built-in map, this type is not concurrency safe
type StringMap struct {
	keys   []string
	values map[string]string
}

// Set sets a key to a value
// If a key already exists it is overwritten
func (m *StringMap) Set(key, value string) {
	if m.values == nil {
		m.keys = append(m.keys, key)
		m.values = map[string]string{key: value}
	} else {
		if _, exists := m.values[key]; !exists {
			m.keys = append(m.keys, key)
		}
		m.values[key] = value
	}
}

// Keys returns the keys in order
func (m StringMap) Keys() []string {
	keys := make([]string, len(m.keys))
	copy(keys, m.keys)

	return keys
}

// Value returns the value for key
func (m StringMap) Value(key string) (string, bool) {
	value, ok := m.values[key]
	return value, ok
}

// Sort sorts the list by value using the provided function
func (m *StringMap) Sort(less func(s, t string) bool) {
	sort.Slice(m.keys, func(i, j int) bool {
		// Use the value for sorting
		return less(m.values[m.keys[i]], m.values[m.keys[j]])
	})
}

// SortKeys sorts the list by key using the provided function
func (m *StringMap) SortKeys(less func(s, t string) bool) {
	sort.Slice(m.keys, func(i, j int) bool {
		return less(m.keys[i], m.keys[j])
	})
}

// MarshalJSON implements json.Marshaler
func (m StringMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("{")
	for i, key := range m.keys {
		var bKey, bVal []byte
		if i > 0 {
			buf.WriteString(",")
		}

		// marshal key
		bKey, _ = json.Marshal(key)
		buf.Write(bKey)
		buf.WriteString(":")

		// marshal value
		bVal, _ = json.Marshal(m.values[key])
		buf.Write(bVal)
	}
	buf.WriteString("}")

	return buf.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler
func (m *StringMap) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))

	// start of object
	if t, err := d.Token(); err != nil {
		return err
	} else if t != json.Delim('{') {
		return errors.New("looking for beginning of object")
	}

	// key/value pairs
	for d.More() {
		tKey, err := d.Token()
		if err != nil {
			return err
		}

		tVal, err := d.Token()
		if err != nil {
			return err
		}
		sVal, ok := tVal.(string)
		if !ok {
			return fmt.Errorf("invalid value type %T", tVal)
		}

		m.Set(tKey.(string), sVal)
	}

	// end of object
	if t, err := d.Token(); t != json.Delim('}') {
		return err
	}

	// end of input
	if _, err := d.Token(); err != io.EOF {
		return errors.New("expected end of JSON input")
	}
	return nil
}

// Len is part of sort.Interface
func (m StringMap) Len() int { return len(m.keys) }

// Less is part of sort.Interface
// Implements same behavior as sort.StringSlice
func (m StringMap) Less(i, j int) bool { return m.values[m.keys[i]] < m.values[m.keys[j]] }

// Swap is part of sort.Interface
func (m StringMap) Swap(i, j int) {
	m.keys[i], m.keys[j] = m.keys[j], m.keys[i]
}
