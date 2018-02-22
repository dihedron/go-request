package request

import "testing"

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
