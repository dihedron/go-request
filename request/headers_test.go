package request

import (
	"testing"
)

func TestValueHeadersProviderSetHeader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &ValueHeadersProvider{
		headers: map[string][]string{},
	}
	provider.AddHeader("key", "value1", "value2", "value3")
	// replace values
	provider.SetHeader("key", "value4", "value5")
	if values, ok := provider.headers["key"]; ok {
		if len(values) != 2 {
			t.Fatalf("invalid headers values: %v", values)
		}
		if values[0] != "value4" {
			t.Fatalf("invalid header value: %v (expected 'value4')", values[0])
		}
		if values[1] != "value5" {
			t.Fatalf("invalid header value: %v (expected 'value5')", values[1])
		}
	}
	// remove values
	provider.SetHeader("key")
	if values, ok := provider.headers["key"]; ok {
		t.Fatalf("invalid headers values: there should be none under 'key', got %v instead", values)
	}
}

func TestValueHeadersProviderGetHeader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &ValueHeadersProvider{
		headers: map[string][]string{},
	}
	provider.AddHeader("key", "value1", "value2", "value3")
	// get the header values back
	values := provider.GetHeader("key")
	if len(values) != 3 {
		t.Fatalf("invalid headers values: %v", values)
	}
	if values[0] != "value1" {
		t.Fatalf("invalid header value: %v (expected 'value1')", values[0])
	}
	if values[1] != "value2" {
		t.Fatalf("invalid header value: %v (expected 'value2')", values[1])
	}
	if values[2] != "value3" {
		t.Fatalf("invalid header value: %v (expected 'value3')", values[2])
	}
	// get another non-existing header
	values = provider.GetHeader("another_key")
	if len(values) != 0 {
		t.Fatalf("invalid headers values: expected none, got %v", values)
	}
	// pass an invalid key
	values = provider.GetHeader("")
	if values != nil {
		t.Fatalf("expected nil return, got %v", values)
	}
}

func TestValueHeadersProviderAddHeader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &ValueHeadersProvider{
		headers: map[string][]string{},
	}
	provider.AddHeader("key", "value1", "value2", "value3")
	if values, ok := provider.headers["key"]; ok {
		if len(values) != 3 {
			t.Fatalf("invalid headers values: %v", values)
		}
		if values[0] != "value1" {
			t.Fatalf("invalid header value: %v (expected 'value1')", values[0])
		}
		if values[1] != "value2" {
			t.Fatalf("invalid header value: %v (expected 'value2')", values[1])
		}
		if values[2] != "value3" {
			t.Fatalf("invalid header value: %v (expected 'value3')", values[2])
		}
	}
	// now add more
	provider.AddHeader("key", "value4", "value5")
	if values, ok := provider.headers["key"]; ok {
		if len(values) != 5 {
			t.Fatalf("invalid headers values: %v", values)
		}
		if values[0] != "value1" {
			t.Fatalf("invalid header value: %v (expected 'value1')", values[0])
		}
		if values[1] != "value2" {
			t.Fatalf("invalid header value: %v (expected 'value2')", values[1])
		}
		if values[2] != "value3" {
			t.Fatalf("invalid header value: %v (expected 'value3')", values[2])
		}
		if values[3] != "value4" {
			t.Fatalf("invalid header value: %v (expected 'value4')", values[3])
		}
		if values[4] != "value5" {
			t.Fatalf("invalid header value: %v (expected 'value5')", values[4])
		}
	}
}

func TestValueHeadersProviderRemoveHeader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	provider := &ValueHeadersProvider{
		headers: map[string][]string{},
	}
	provider.AddHeader("key", "value1", "value2", "value3")

	provider.RemoveHeader("key")
	if values, ok := provider.headers["key"]; ok {
		t.Fatalf("header 'key' should not be there: got %v", values)
	}
}

func TestStructHeadersProviderSetSource(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("recovered from correct panic due to '%v'", r)
		}
	}()

	provider := &StructHeadersProvider{}

	source := struct {
		Header1 string `header:"X-Header-1,omitempty"`
		Header2 string `header:"X-Header-2,omitempty"`
	}{}

	provider.SetSource(source)

	provider.SetSource(&source)

	provider.SetSource("no_struct")

	t.Fatal("should not arrive here")
}

func TestStructHeadersProviderGetHeaders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Logf("correctly recovered: %v", r)
		}
	}()

	provider := &StructHeadersProvider{}

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

	provider.SetSource(source)
	headers := provider.GetHeaders()
	for key, values := range headers {
		// t.Logf("%v => %v", key, values)
		switch key {
		case "X-Header-1":
			if len(values) != 2 {
				t.Fatalf("invalid number of header values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-1a" && value != "header-value-1b" {
					t.Fatalf("invalid header value: %v", value)
				}
			}
		case "X-Header-2":
			if len(values) != 2 {
				t.Fatalf("invalid number of header values, expected 2, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-2a" && value != "header-value-2b" {
					t.Fatalf("invalid header value: %v", value)
				}
			}
		case "X-Header-3":
			if len(values) != 2 {
				t.Fatalf("invalid number of header values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-3a" && value != "header-value-3b" {
					t.Fatalf("invalid header value: %v", value)
				}
			}
		case "X-Header-4":
			if len(values) != 1 {
				t.Fatalf("invalid number of header values, expected 1, got %d", len(values))
			}
			for _, value := range values {
				if value != "header-value-4" {
					t.Fatalf("invalid header value: %v", value)
				}
			}
		}
	}
}
