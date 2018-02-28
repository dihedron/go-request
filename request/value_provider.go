// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"fmt"
	"reflect"

	"github.com/fatih/structs"
)

// MapProvider is a simplistic implementation of a Provider, based on a map where
// key/value pairs are stored.
type MapProvider struct {
	values map[string][]interface{}
}

// Provide returns the map of currently registered keys and values.
func (p MapProvider) Provide() map[string][]string {
	result := map[string][]string{}
	for key := range p.values {
		result[key] = p.GetAsStrings(key)
	}
	return result
}

// Set replaces any existing value of the given key with the values provided; if
// no value is provided, this method is equivalent to dropping the key.
func (p MapProvider) Set(key string, values ...interface{}) MapProvider {
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
func (p MapProvider) Get(key string) []interface{} {
	if key != "" {
		return p.values[key]
	}
	return nil
}

// GetAsStrings returns the values associated with the given key, converted
// strings using the fmt package.
func (p MapProvider) GetAsStrings(key string) []string {
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
func (p MapProvider) Add(key string, values ...string) MapProvider {
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
func (p MapProvider) Remove(key string) MapProvider {
	if key != "" {
		delete(p.values, key)
	}
	return p
}

// StructProvider extracts key/value pairs from a struct by inspecting its
// fields and looking for those annotated with the given tag; the struct is
// scanned recursively, so both embedded and nested structs are visited.
type StructProvider struct {
	// tag is the tag to look for when scanning the struct.
	Tag string
	// source is the struct used to retrieve keys and values.
	Source interface{}
}

// Provide extracts the headers from the source struct and returns them.
func (p StructProvider) Provide() map[string][]interface{} {
	var source interface{}
	switch reflect.ValueOf(p.Source).Kind() {
	case reflect.Struct:
		source = p.Source
	case reflect.Ptr:
		if reflect.ValueOf(p.Source).Elem().Kind() == reflect.Struct {
			source = reflect.ValueOf(p.Source).Elem().Interface()
		} else {
			panic("only structs can be passed as providers")
		}
	default:
		panic("only structs can be passed as providers")
	}

	if p.Tag == "" {
		panic("a valid tag must be provided")
	}

	return scan(p.Tag, source)
}

// scan is the actual workhorse method: it scans the source struct for tagged
// headers and extracts their values; if any embedded or child struct is
// encountered, it is scanned for values.
func scan(key string, source interface{}) map[string][]interface{} {
	result := map[string][]interface{}{}
	for _, field := range structs.Fields(source) {
		if field.IsEmbedded() || field.Kind() == reflect.Struct {
			for k, v := range scan(key, field.Value()) {
				if values, ok := result[k]; ok {
					result[k] = append(values, v...)
				} else {
					result[k] = v
				}
			}
		} else {
			tag := field.Tag(key)
			if tag != "" {
				value := field.Value()
				if values, ok := result[tag]; ok {
					result[tag] = append(values, value)
				} else {
					result[tag] = []interface{}{value}
				}
			}
		}
	}
	return result
}
