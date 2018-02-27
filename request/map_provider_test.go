// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"fmt"
	"testing"
)

func TestMapProcviderProvide(t *testing.T) {
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
