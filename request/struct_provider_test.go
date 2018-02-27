// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import "testing"

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
