// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"net/http"
	"reflect"

	"github.com/fatih/structs"
)

// Headers represents the HTTP headers that will go along with the request.
type Headers http.Header

// HeadersProvider is the common interface to all entities capable of generating
// or providinig headers for the request.
type HeadersProvider interface {
	// ProvideHeaders returns a Headers map; the keys and values in the returned
	// map will be set as request headers.
	GetHeaders() Headers
}

// ValueHeadersProvider is a simplistic implementation of a HeadersProvider;
// is stores key/value pairs and then returns them.
type ValueHeadersProvider struct {
	// headers is the set of key/value pairs.
	headers Headers
}

// SetHeader replaces any existing value of the given header with the values
// provided; if no value is provided, it is equivalent to dropping the header.
func (p *ValueHeadersProvider) SetHeader(key string, values ...string) *ValueHeadersProvider {
	if key != "" {
		if len(values) > 0 {
			p.headers[key] = values
		} else {
			delete(p.headers, key)
		}
	}
	return p
}

// GetHeader returns the values associated with the given key.
func (p *ValueHeadersProvider) GetHeader(key string) []string {
	if key != "" {
		return p.headers[key]
	}
	return nil
}

// AddHeader adds the provided values to the given header; if none exists already
// with the given key, it is created.
func (p *ValueHeadersProvider) AddHeader(key string, values ...string) *ValueHeadersProvider {
	if key != "" {
		if existing, ok := p.headers[key]; ok {
			p.headers[key] = append(existing, values...)
		} else {
			p.headers[key] = values
		}
	}
	return p
}

// RemoveHeader removes the header with the given key from the map, if it exists;
// if it doesn't, it won't do anything.
func (p *ValueHeadersProvider) RemoveHeader(key string) *ValueHeadersProvider {
	if key != "" {
		delete(p.headers, key)
	}
	return p
}

// StructHeadersProvider extracts headers value from a struct by inspecting its
// fields and looking for thos annotated with the 'header' tag; note: the struct
// is not scanned recursively, so only first level fields are used.
type StructHeadersProvider struct {
	// source if the struct used to retrieve header values.
	source interface{}
}

// SetSource sets the
func (p *StructHeadersProvider) SetSource(source interface{}) *StructHeadersProvider {
	switch reflect.ValueOf(source).Kind() {
	case reflect.Struct:
		p.source = source
		return p
	case reflect.Ptr:
		if reflect.ValueOf(source).Elem().Kind() == reflect.Struct {
			p.source = reflect.ValueOf(source).Elem()
			return p
		}
	}
	panic("only structs can be passed as headers provider")
}

// GetHeaders extracts the headers from the source struct and returns them.
func (p *StructHeadersProvider) GetHeaders() Headers {
	return getHeaders(p.source)
}

// getHeaders is the actual workhorse method: it scans the source struct for
// tagged headers and extracts their values; if any embedded or child struct is
// encountered, it is scanned for values.
func getHeaders(source interface{}) Headers {
	headers := map[string][]string{}
	for _, field := range structs.Fields(source) {
		if field.IsEmbedded() || field.Kind() == reflect.Struct {
			for k, v := range getHeaders(field.Value()) {
				if values, ok := headers[k]; ok {
					headers[k] = append(values, v...)
				} else {
					headers[k] = v
				}
			}
		} else {
			tag := field.Tag("header")
			if tag != "" {
				value := field.Value().(string)
				if values, ok := headers[tag]; ok {
					headers[tag] = append(values, value)
				} else {
					headers[tag] = []string{value}
				}
			}
		}
	}
	return headers
}
