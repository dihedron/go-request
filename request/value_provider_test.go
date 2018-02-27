// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"fmt"
	"testing"
)

func TestMapProviderProvide(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key1", "value1.1")
	provider.Add("key2", "value2.1", "value2.2")
	provider.Add("key3", "value3.1", "value3.2", "value3.3")
	provider.Add("key4", "value4.1", "value4.2", "value4.3", "value4.4")
	provider.Add("key5", "value5.1", "value5.2", "value5.3", "value5.4", "value5.5")
	provider.Add("key6", "value6.1", "value6.2", "value6.3", "value6.4", "value6.5", "value6.6")

	values := provider.Provide()
	if len(values) != 6 {
		t.Fatalf("invalid values map size: expected 3, got %v", len(values))
	}
	for i := 1; i <= len(values); i++ {
		key := fmt.Sprintf("key%d", i)
		list := values[key]
		if len(list) != i {
			t.Fatalf("invalid number of values for %s: expected %d, got %d", key, i, len(list))
		}
		for j, actual := range list {
			expected := fmt.Sprintf("value%d.%d", i, j+1)
			if actual != expected {
				t.Fatalf("invalid value: expected %s, got %s", expected, actual)
			}
		}
	}

}

func TestMapProviderSet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key", "value1", "value2", "value3")
	// replace values
	provider.Set("key", "value4", "value5")
	if values, ok := provider.values["key"]; ok {
		if len(values) != 2 {
			t.Fatalf("invalid values: %v", values)
		}
		if values[0] != "value4" {
			t.Fatalf("invalid value: %v (expected 'value4')", values[0])
		}
		if values[1] != "value5" {
			t.Fatalf("invalid value: %v (expected 'value5')", values[1])
		}
	}
	// remove values
	provider.Set("key")
	if values, ok := provider.values["key"]; ok {
		t.Fatalf("invalid values: there should be none under 'key', got %v instead", values)
	}
}

func TestMapProviderGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key", "value1", "value2", "value3")
	// get the values back
	values := provider.Get("key")
	if len(values) != 3 {
		t.Fatalf("invalid values: %v", values)
	}
	if values[0] != "value1" {
		t.Fatalf("invalid value: %v (expected 'value1')", values[0])
	}
	if values[1] != "value2" {
		t.Fatalf("invalid value: %v (expected 'value2')", values[1])
	}
	if values[2] != "value3" {
		t.Fatalf("invalid value: %v (expected 'value3')", values[2])
	}
	// get another non-existing header
	values = provider.Get("another_key")
	if len(values) != 0 {
		t.Fatalf("invalid values: expected none, got %v", values)
	}
	// pass an invalid key
	values = provider.Get("")
	if values != nil {
		t.Fatalf("expected nil return, got %v", values)
	}
}

func TestMapProviderGetAsStrings(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key", "value1", "value2", "value3")
	// get the values back
	values := provider.GetAsStrings("key")
	if len(values) != 3 {
		t.Fatalf("invalid values: %v", values)
	}
	if values[0] != "value1" {
		t.Fatalf("invalid value: %v (expected 'value1')", values[0])
	}
	if values[1] != "value2" {
		t.Fatalf("invalid value: %v (expected 'value2')", values[1])
	}
	if values[2] != "value3" {
		t.Fatalf("invalid value: %v (expected 'value3')", values[2])
	}
	// get another non-existing header
	values = provider.GetAsStrings("another_key")
	if len(values) != 0 {
		t.Fatalf("invalid values: expected none, got %v", values)
	}
	// pass an invalid key
	values = provider.GetAsStrings("")
	if values != nil {
		t.Fatalf("expected nil return, got %v", values)
	}
}

func TestMapProviderAdd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key", "value1", "value2", "value3")
	if values, ok := provider.values["key"]; ok {
		if len(values) != 3 {
			t.Fatalf("invalid values: %v", values)
		}
		if values[0] != "value1" {
			t.Fatalf("invalid value: %v (expected 'value1')", values[0])
		}
		if values[1] != "value2" {
			t.Fatalf("invalid value: %v (expected 'value2')", values[1])
		}
		if values[2] != "value3" {
			t.Fatalf("invalid value: %v (expected 'value3')", values[2])
		}
	}
	// now add more
	provider.Add("key", "value4", "value5")
	if values, ok := provider.values["key"]; ok {
		if len(values) != 5 {
			t.Fatalf("invalid values: %v", values)
		}
		if values[0] != "value1" {
			t.Fatalf("invalid value: %v (expected 'value1')", values[0])
		}
		if values[1] != "value2" {
			t.Fatalf("invalid value: %v (expected 'value2')", values[1])
		}
		if values[2] != "value3" {
			t.Fatalf("invalid value: %v (expected 'value3')", values[2])
		}
		if values[3] != "value4" {
			t.Fatalf("invalid value: %v (expected 'value4')", values[3])
		}
		if values[4] != "value5" {
			t.Fatalf("invalid value: %v (expected 'value5')", values[4])
		}
	}
}

func TestMapProviderRemove(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &MapProvider{
		values: map[string][]interface{}{},
	}
	provider.Add("key", "value1", "value2", "value3")

	provider.Remove("key")
	if values, ok := provider.values["key"]; ok {
		t.Fatalf("'key' should not be there: got %v", values)
	}
}

func TestStructProviderProvideStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	type Embedded struct {
		NonHeader string
		Header1   string `header:"X-Header-1"`
		Header2   string `header:"X-Header-2"`
		Header3a  string `header:"X-Header-3"`
		Header3b  string `header:"X-Header-3"`
	}

	type Nested struct {
		NonHeader string
		Header4   string `header:"X-Header-4"`
	}

	source := struct {
		Header1 string `header:"X-Header-1"`
		Header2 string `header:"X-Header-2"`
		Embedded
		Struct Nested
	}{
		Header1: "header-value-1a",
		Header2: "header-value-2a",
		Embedded: Embedded{
			NonHeader: "non-header-value",
			Header1:   "header-value-1b",
			Header2:   "header-value-2b",
			Header3a:  "header-value-3a",
			Header3b:  "header-value-3b",
		},
		Struct: Nested{
			NonHeader: "non-header-value",
			Header4:   "header-value-4",
		},
	}

	// this call is expected to succeed
	provider := &StructProvider{
		Source: source,
		Tag:    "header",
	}
	result := provider.Provide()
	for key, values := range result {
		switch key {
		case "X-Header-1":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-1a" && value != "header-value-1b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-2":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-2a" && value != "header-value-2b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-3":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-3a" && value != "header-value-3b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-4":
			if len(values) != 1 {
				t.Fatalf("invalid number of values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-4" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		}
	}
}

func TestStructProviderProvidePtr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	type Embedded struct {
		NonHeader string
		Header1   string `header:"X-Header-1"`
		Header2   string `header:"X-Header-2"`
		Header3a  string `header:"X-Header-3"`
		Header3b  string `header:"X-Header-3"`
	}

	type Nested struct {
		NonHeader string
		Header4   string `header:"X-Header-4"`
	}

	source := &struct {
		Header1 string `header:"X-Header-1"`
		Header2 string `header:"X-Header-2"`
		Embedded
		Struct Nested
	}{
		Header1: "header-value-1a",
		Header2: "header-value-2a",
		Embedded: Embedded{
			NonHeader: "non-header-value",
			Header1:   "header-value-1b",
			Header2:   "header-value-2b",
			Header3a:  "header-value-3a",
			Header3b:  "header-value-3b",
		},
		Struct: Nested{
			NonHeader: "non-header-value",
			Header4:   "header-value-4",
		},
	}

	// this call is expected to succeed
	provider := &StructProvider{
		Source: source,
		Tag:    "header",
	}
	result := provider.Provide()
	for key, values := range result {
		switch key {
		case "X-Header-1":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-1a" && value != "header-value-1b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-2":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-2a" && value != "header-value-2b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-3":
			if len(values) != 2 {
				t.Fatalf("invalid number of values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-3a" && value != "header-value-3b" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		case "X-Header-4":
			if len(values) != 1 {
				t.Fatalf("invalid number of values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-4" {
					t.Fatalf("invalid value: %v", value)
				}
			}
		}
	}

}

func TestStructProviderProvideNoStruct(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer handler("only structs can be passed as providers", t)

	// this call is expected to fail due to a non-struct source
	provider := &StructProvider{
		Tag:    "header",
		Source: "no_struct",
	}
	provider.Provide()
}

func TestStructProviderProvideNoStructPtr(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer handler("only structs can be passed as providers", t)

	// this call is expected to fail due to a non-struct source
	s := "no_struct"
	provider := &StructProvider{
		Tag:    "header",
		Source: &s,
	}
	provider.Provide()
}

func TestStructProviderProvideNoTag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	type Embedded struct {
		NonHeader string
		Header1   string `header:"X-Header-1"`
		Header2   string `header:"X-Header-2"`
		Header3a  string `header:"X-Header-3"`
		Header3b  string `header:"X-Header-3"`
	}

	type Nested struct {
		NonHeader string
		Header4   string `header:"X-Header-4"`
	}

	source := struct {
		Header1 string `header:"X-Header-1"`
		Header2 string `header:"X-Header-2"`
		Embedded
		Struct Nested
	}{
		Header1: "header-value-1a",
		Header2: "header-value-2a",
		Embedded: Embedded{
			NonHeader: "non-header-value",
			Header1:   "header-value-1b",
			Header2:   "header-value-2b",
			Header3a:  "header-value-3a",
			Header3b:  "header-value-3b",
		},
		Struct: Nested{
			NonHeader: "non-header-value",
			Header4:   "header-value-4",
		},
	}

	defer handler("a valid tag must be provided", t)

	// this call is expected to succeed
	provider := &StructProvider{
		Source: source,
	}
	provider.Provide()
}

func handler(message string, t *testing.T) {
	if r := recover(); r != nil {
		if r == message {
			t.Logf("correctly recovered: %v", r)
		} else {
			t.Fatalf("unxpected panic: %v", r)
		}
	}
}
