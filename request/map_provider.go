// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import "fmt"

// MapProvider is a simplistic implementation of a Provider, based on a map where
// key/value pairs are stored.
type MapProvider struct {
	values map[string][]interface{}
}

// Provide returns the map of currently registered keys and values.
func (p *MapProvider) Provide() map[string][]string {
	result := map[string][]string{}
	for key := range p.values {
		result[key] = p.GetAsStrings(key)
	}
	return result
}

// Set replaces any existing value of the given key with the values provided; if
// no value is provided, this method is equivalent to dropping the key.
func (p *MapProvider) Set(key string, values ...interface{}) *MapProvider {
	if key != "" {
		if len(values) > 0 {
			p.values[key] = values
		} else {
			delete(p.values, key)
		}
	}
	return p
}

// Get returns the values associated with the given key.
func (p *MapProvider) Get(key string) []interface{} {
	if key != "" {
		return p.values[key]
	}
	return nil
}

// GetAsStrings returns the values associated with the given key, converted
// strings using the fmt package.
func (p *MapProvider) GetAsStrings(key string) []string {
	if key != "" {
		values := []string{}
		for _, value := range p.values[key] {
			values = append(values, fmt.Sprintf("%v", value))
		}
		return values
	}
	return nil
}

// Add adds the provided values to the given key; if no value exists already
// under the given key, the entry is created.
func (p *MapProvider) Add(key string, values ...string) *MapProvider {
	if key != "" {
		temp := []interface{}{}
		for _, value := range values {
			temp = append(temp, value)
		}
		if existing, ok := p.values[key]; ok {
			p.values[key] = append(existing, temp...)
		} else {
			p.values[key] = temp
		}
	}
	return p
}

// Remove removes the entry with the given key from the map, if it exists; if
// it doesn't, it doesn't do anything.
func (p *MapProvider) Remove(key string) *MapProvider {
	if key != "" {
		delete(p.values, key)
	}
	return p
}
