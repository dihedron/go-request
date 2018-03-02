// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request2

import (
	"net"
	"net/http"
	"testing"
	"time"
)

func getClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
}

func TestNewFactory(t *testing.T) {
	f := New("https://www.example.com")
	if f.method != http.MethodGet {
		t.Fatalf("invalid method: expected get, got %v", f.method)
	}
	if f.url != "https://www.example.com" {
		t.Fatalf("invalid url: expected \"https://www.example.com\", got %v", f.url)
	}
	if len(f.parameters) != 0 {
		t.Fatalf("invalid query parameters: expected empty, got %d entries", len(f.parameters))
	}
	if len(f.headers) != 0 {
		t.Fatalf("invalid headers: expected empty, got %d entries", len(f.headers))
	}

	f = New("")
	if f.url != "" {
		t.Fatalf("invalid url: expected empty, got %v", f.url)
	}
}

func TestNewSubFactory(t *testing.T) {
	f := New("https://www.example.com").New("", "")
	if f.method != http.MethodGet {
		t.Fatalf("invalid method: expected GET, got %v", f.method)
	}
	if f.url != "https://www.example.com" {
		t.Fatalf("invalid url: expected \"https://www.example.com\", got %v", f.url)
	}
	if len(f.parameters) != 0 {
		t.Fatalf("invalid query parameters: expected empty, got %d entries", len(f.parameters))
	}
	if len(f.headers) != 0 {
		t.Fatalf("invalid headers: expected empty, got %d entries", len(f.headers))
	}

	f = New("https://www.example.com").New(http.MethodPost, "https://www.example.com/sub")
	if f.method != http.MethodPost {
		t.Fatalf("invalid method: expected POST, got %v", f.method)
	}
	if f.url != "https://www.example.com/sub" {
		t.Fatalf("invalid url: expected \"https://www.example.com/sub\", got %v", f.url)
	}
	if len(f.parameters) != 0 {
		t.Fatalf("invalid query parameters: expected empty, got %d entries", len(f.parameters))
	}
	if len(f.headers) != 0 {
		t.Fatalf("invalid headers: expected empty, got %d entries", len(f.headers))
	}
}

func TestBase(t *testing.T) {
	f := New("")
	if f.url != "" {
		t.Fatalf("invalid url: expected \"\", got %v", f.url)
	}
	f.Base("https://www.example.com")
	if f.url != "https://www.example.com" {
		t.Fatalf("invalid url: expected \"https://www.example.com\", got %v", f.url)
	}
}

func TestPath(t *testing.T) {
	f := New("https://www.example.com").Path("/api/v2")
	if f.url != "https://www.example.com/api/v2" {
		t.Fatalf("invalid url: expected \"https://www.example.com/api/v2\", got %v", f.url)
	}

	f = New("https://www.example.com/").Path("/api/v2")
	if f.url != "https://www.example.com/api/v2" {
		t.Fatalf("invalid url: expected \"https://www.example.com/api/v2\", got %v", f.url)
	}

	f = New("https://www.example.com/test").Path("../api/v2")
	if f.url != "https://www.example.com/api/v2" {
		t.Fatalf("invalid url: expected \"https://www.example.com/api/v2\", got %v", f.url)
	}

	f = New("https://www.example.com/test/").Path("../api/v2")
	if f.url != "https://www.example.com/api/v2" {
		t.Fatalf("invalid url: expected \"https://www.example.com/api/v2\", got %v", f.url)
	}

	f = New("https://www.example.com/test/").Path("https://www.example.com/api/v2")
	if f.url != "https://www.example.com/api/v2" {
		t.Fatalf("invalid url: expected \"https://www.example.com/api/v2\", got %v", f.url)
	}
	// I couldn't find a way to produce a URL that cannot be parsed somehow...
}

func TestUserAgent(t *testing.T) {
	expected := "MyCrawler/1.0"
	f := New("").UserAgent(expected)
	if len(f.headers["User-Agent"]) != 1 {
		t.Fatalf("user agent header should have 1 value, got %d", len(f.headers["User-Agent"]))
	}
	if value := f.headers["User-Agent"][0]; value != expected {
		t.Fatalf("invalid user agent header: expected %q, got %q", expected, value)
	}
}

func TestContentType(t *testing.T) {
	expected := "text/xml"
	f := New("").ContentType(expected)
	if len(f.headers["Content-Type"]) != 1 {
		t.Fatalf("content type should have 1 value, got %d", len(f.headers["Content-Type"]))
	}
	if value := f.headers["Content-Type"][0]; value != expected {
		t.Fatalf("invalid content type header: expected %q, got %q", expected, value)
	}
}

func TestMethod(t *testing.T) {
	f := New("")
	tests := []struct {
		method   func() *Factory
		expected string
	}{
		{f.Get, "GET"},
		{f.Post, "POST"},
		{f.Put, "PUT"},
		{f.Delete, "DELETE"},
		{f.Patch, "PATCH"},
		{f.Options, "OPTIONS"},
		{f.Trace, "TRACE"},
		{f.Connect, "CONNECT"},
		{f.Head, "HEAD"},
	}

	for _, test := range tests {
		test.method()
		if f.method != test.expected {
			t.Fatalf("invalid method: expected %s, got %s", test.expected, f.method)
		}
	}
}

func TestOp(t *testing.T) {
	f := New("")
	tests := []struct {
		method   func() *Factory
		expected operation
	}{
		{f.Add, add},
		{f.Set, set},
		{f.Del, del},
		{f.Remove, rem},
	}

	for _, test := range tests {
		test.method()
		if f.op != test.expected {
			t.Fatalf("invalid operation: expected %d, got %d", test.expected, f.op)
		}
	}
}
