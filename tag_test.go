// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"testing"

	"github.com/fatih/structs"
)

func TestTag(t *testing.T) {
	s := struct {
		A string `mytag:"myvalue"`
		B string `mytag:"myvalue,omitempty"`
		C string `mytag:"-"`
		D string `mytag:"-,omitempty"`
		E string
	}{
		A: "a",
		B: "b",
		C: "c",
		D: "d",
		E: "e",
	}
	for _, field := range structs.Fields(s) {
		tag := NewTag(field.Tag("mytag"))
		switch field.Name() {
		case "A":
			if tag.Name() != "myvalue" {
				t.Fatalf("field A must have \"myvalue\"")
			}
			if tag.IsIgnore() {
				t.Fatalf("field A is not ignorable")
			}
			if tag.IsOmitEmpty() {
				t.Fatalf("field A is not omitempty")
			}
		case "B":
			if tag.Name() != "myvalue" {
				t.Fatalf("field B must have \"myvalue\"")
			}
			if tag.IsIgnore() {
				t.Fatalf("field B is not ignorable")
			}
			if !tag.IsOmitEmpty() {
				t.Fatalf("field B is omitempty")
			}
		case "C":
			if tag.Name() != "-" {
				t.Fatalf("field C must have \"-\"")
			}
			if !tag.IsIgnore() {
				t.Fatalf("field C is ignorable")
			}
			if tag.IsOmitEmpty() {
				t.Fatalf("field C is not omitempty")
			}
		case "D":
			if tag.Name() != "-" {
				t.Fatalf("field D must have \"-\"")
			}
			if !tag.IsIgnore() {
				t.Fatalf("field D is ignorable")
			}
			if !tag.IsOmitEmpty() {
				t.Fatalf("field D is omitempty")
			}
		case "E":
			if tag.Name() != "" {
				t.Fatalf("field E must be not valid (\"\")")
			}
			if !tag.IsMissing() {
				t.Fatalf("field E must have a missing tag")
			}
		}
	}
}
