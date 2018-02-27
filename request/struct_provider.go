// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"reflect"

	"github.com/fatih/structs"
)

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
func (p *StructProvider) Provide() map[string][]interface{} {
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
