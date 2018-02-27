// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

const expectedPersonJSON string = `{
  "name": "John",
  "surname": "Doe",
  "address": {
    "street": "Madison ave.",
    "number": 123,
    "zip": "00123"
  }
}`

func TestJSONProviderProvideStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	type Address struct {
		Street  string `json:"street,omitempty"`
		No      int    `json:"number,omitempty"`
		ZIPCode string `json:"zip,omitempy"`
	}

	type Person struct {
		Name    string  `json:"name,omitempty"`
		Surname string  `json:"surname,omitempty"`
		Address Address `json:"address,omitempty"`
	}

	provider := &JSONProvider{
		Entity: Person{
			Name:    "John",
			Surname: "Doe",
			Address: Address{
				Street:  "Madison ave.",
				No:      123,
				ZIPCode: "00123",
			},
		},
	}
	entity, err := provider.Provide()

	if err != nil {
		t.Fatalf("error encoding to JSON: %v", err)
	}

	if entity.ContentType != "application/json" {
		t.Fatalf("invalid content type returned: %v", entity.ContentType)
	}

	data, _ := ioutil.ReadAll(entity.Reader)
	person := &Person{}
	json.Unmarshal(data, person)
	actual, _ := json.MarshalIndent(person, "", "  ")

	if string(actual) != expectedPersonJSON {
		t.Fatalf("invalid JSON roundtrip: expected:\n%s\nactual:\n%s", expectedPersonJSON, string(actual))
	}
}

func TestJSONProviderProvidePtr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	type Address struct {
		Street  string `json:"street,omitempty"`
		No      int    `json:"number,omitempty"`
		ZIPCode string `json:"zip,omitempy"`
	}

	type Person struct {
		Name    string  `json:"name,omitempty"`
		Surname string  `json:"surname,omitempty"`
		Address Address `json:"address,omitempty"`
	}

	provider := &JSONProvider{
		Entity: &Person{
			Name:    "John",
			Surname: "Doe",
			Address: Address{
				Street:  "Madison ave.",
				No:      123,
				ZIPCode: "00123",
			},
		},
	}
	entity, err := provider.Provide()

	if err != nil {
		t.Fatalf("error encoding to JSON: %v", err)
	}

	if entity.ContentType != "application/json" {
		t.Fatalf("invalid content type returned: %v", entity.ContentType)
	}

	data, _ := ioutil.ReadAll(entity.Reader)
	person := &Person{}
	json.Unmarshal(data, person)
	actual, _ := json.MarshalIndent(person, "", "  ")

	if string(actual) != expectedPersonJSON {
		t.Fatalf("invalid JSON roundtrip: expected:\n%s\nactual:\n%s", expectedPersonJSON, string(actual))
	}
}

func TestJSONProviderProvideNoStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer handler("only structs can be passed as providers", t)

	provider := &JSONProvider{
		Entity: "no_struct",
	}
	provider.Provide()
}

func TestJSONProviderProvideNoPtr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer handler("only structs can be passed as providers", t)

	s := "no_struct"
	provider := &JSONProvider{
		Entity: &s,
	}
	provider.Provide()
}
