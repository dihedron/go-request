// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"bytes"
	"encoding/json"
	"reflect"
)

// JSONProvider is a specialisation of EntityProvider for JSON payloads.
type JSONProvider struct {
	Entity interface{}
}

// Provide converts a struct or a struct pointer to a JSON document, and returns
// an entity containing the JSON in serialiside form.
func (p *JSONProvider) Provide() (*Entity, error) {
	var source interface{}
	switch reflect.ValueOf(p.Entity).Kind() {
	case reflect.Struct:
		source = p.Entity
	case reflect.Ptr:
		if reflect.ValueOf(p.Entity).Elem().Kind() == reflect.Struct {
			source = reflect.ValueOf(p.Entity).Elem().Interface()
		} else {
			panic("only structs can be passed as providers")
		}
	default:
		panic("only structs can be passed as providers")
	}

	data, err := json.Marshal(source)
	if err != nil {
		return nil, err
	}

	return &Entity{
		ContentType: "application/json",
		Reader:      bytes.NewReader(data),
	}, nil
}
