// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"reflect"
)

// JSONProvider is a specialisation of EntityProvider for JSON payloads.
type JSONProvider struct {
	// ContentType is an optional field used when the default "application/xml"
	// content type will not do.
	ContentType string
	// Entity is the entity to be converted to JSON.
	Entity interface{}
}

// Provide converts a struct or a struct pointer to a JSON document, and returns
// an entity containing the JSON in serialised form.
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

	entity := &Entity{
		ContentType: "application/json",
		Reader:      bytes.NewReader(data),
	}
	if p.ContentType != "" {
		entity.ContentType = p.ContentType
	}
	return entity, nil
}

// XMLProvider is a specialisation of EntityProvider for XML payloads.
type XMLProvider struct {
	// ContentType is an optional field used when the default "application/xml"
	// content type will not do.
	ContentType string
	// Entity is the entity to be converted to XML.
	Entity interface{}
}

// Provide converts a struct or a struct pointer to an XML document, and returns
// an entity containing the XML in serialised form.
func (p *XMLProvider) Provide() (*Entity, error) {
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

	data, err := xml.Marshal(source)
	if err != nil {
		return nil, err
	}

	entity := &Entity{
		ContentType: "application/xml",
		Reader:      bytes.NewReader(data),
	}
	if p.ContentType != "" {
		entity.ContentType = p.ContentType
	}
	return entity, nil
}
