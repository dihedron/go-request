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
	f := NewFactory("https://www.example.com")
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

	f = NewFactory("")
	if f.url != "" {
		t.Fatalf("invalid url: expected empty, got %v", f.url)
	}
}

func TestNewSubFactory(t *testing.T) {
	f := NewFactory("https://www.example.com").New("", "")
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

	f = NewFactory("https://www.example.com").New(http.MethodPost, "https://www.example.com/sub")
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
