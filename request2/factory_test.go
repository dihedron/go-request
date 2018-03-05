// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request2

import (
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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

func TestAddQueryParameter(t *testing.T) {
	f := New("").Add().QueryParameter("key", "value1", "value2")
	if len(f.parameters) != 1 {
		t.Fatalf("error adding query parameters: expected 1, got %d", len(f.parameters))
	}
	if len(f.parameters["key"]) != 2 {
		t.Fatalf("error adding query parameters: expected 2, got %d", len(f.parameters["key"]))
	}
	if f.parameters["key"][0] != "value1" {
		t.Fatalf("error adding query parameters: expected \"value1\", got %q", f.parameters["key"][0])
	}
	if f.parameters["key"][1] != "value2" {
		t.Fatalf("error adding query parameters: expected \"value2\", got %q", f.parameters["key"][1])
	}
	// now add one more
	f.QueryParameter("key", "value3")
	if len(f.parameters) != 1 {
		t.Fatalf("error adding query parameters: expected 1, got %d", len(f.parameters))
	}
	if len(f.parameters["key"]) != 3 {
		t.Fatalf("error adding query parameters: expected 3, got %d", len(f.parameters["key"]))
	}
	if f.parameters["key"][0] != "value1" {
		t.Fatalf("error adding query parameters: expected \"value1\", got %q", f.parameters["key"][0])
	}
	if f.parameters["key"][1] != "value2" {
		t.Fatalf("error adding query parameters: expected \"value2\", got %q", f.parameters["key"][1])
	}
	if f.parameters["key"][2] != "value3" {
		t.Fatalf("error adding query parameters: expected \"value3\", got %q", f.parameters["key"][2])
	}
}

func TestSetQueryParameter(t *testing.T) {
	f := New("").Add().QueryParameter("key", "value1", "value2")
	f.Set().QueryParameter("key", "value0", "value1", "value2", "value3")
	if len(f.parameters) != 1 {
		t.Fatalf("error setting query parameters: expected 1, got %d", len(f.parameters))
	}
	if len(f.parameters["key"]) != 4 {
		t.Fatalf("error setting query parameters: expected 4, got %d", len(f.parameters["key"]))
	}
	if f.parameters["key"][0] != "value0" {
		t.Fatalf("error setting query parameters: expected \"value0\", got %q", f.parameters["key"][0])
	}
	if f.parameters["key"][1] != "value1" {
		t.Fatalf("error setting query parameters: expected \"value1\", got %q", f.parameters["key"][1])
	}
	if f.parameters["key"][2] != "value2" {
		t.Fatalf("error setting query parameters: expected \"value2\", got %q", f.parameters["key"][2])
	}
	if f.parameters["key"][3] != "value3" {
		t.Fatalf("error setting query parameters: expected \"value3\", got %q", f.parameters["key"][3])
	}
}

func TestDelQueryParameter(t *testing.T) {
	f := New("").Add().QueryParameter("key", "value1", "value2")
	f.Del().QueryParameter("key")
	if len(f.parameters) != 0 {
		t.Fatalf("error deleting query parameters: expected 0, got %d", len(f.parameters))
	}
}

func TestRemQueryParameter(t *testing.T) {
	f := New("").
		Add().
		QueryParameter("key1", "value1", "value2").
		QueryParameter("key2", "value1", "value2").
		QueryParameter("key3", "value1", "value2").
		QueryParameter("another_key", "value1", "value2")
	f.Remove().QueryParameter("^key\\d$")
	if len(f.parameters) != 1 {
		t.Fatalf("error removing multiple query parameters: expected 1, got %d", len(f.parameters))
	}
}

func TestQueryParameterFrom(t *testing.T) {

	type Nested struct {
		Query7 string `parameter:"query7"`
		Query1 string `parameter:"query1"`
	}

	type Embedded struct {
		Query5 string `parameter:"query5"`
		Query6 string `parameter:"query6"`
	}

	type Struct struct {
		Query1  string  `parameter:"query1"`
		Query2  *string `parameter:"query2"`
		Query3a bool    `parameter:"query3"`
		Query3b bool    `parameter:"query3"`
		Query4  *bool   `parameter:"query4"`
		Embedded
		Nested *Nested
	}

	s := "value2"
	b := true
	testStruct := Struct{
		Query1:  "value1a",
		Query2:  &s,
		Query3a: true,
		Query3b: false,
		Query4:  &b,
		Embedded: Embedded{
			Query5: "value5",
			Query6: "value6",
		},
		Nested: &Nested{
			Query7: "value7",
			Query1: "value1b",
		},
	}

	testMap := map[string][]string{
		"query1": []string{"value1a", "value1b"},
		"query2": []string{"value2"},
		"query3": []string{"true", "false"},
		"query4": []string{"true"},
		"query5": []string{"value5"},
		"query6": []string{"value6"},
		"query7": []string{"value7"},
	}

	factories := []*Factory{
		New("").Add().QueryParametersFrom(testStruct),
		New("").Add().QueryParametersFrom(&testStruct),
		New("").Add().QueryParametersFrom(testMap),
		New("").Add().QueryParametersFrom(&testMap),
	}

	for _, f := range factories {
		if len(f.parameters) != 7 {
			t.Fatalf("error adding query parameters from struct: expected 7, got %d", len(f.parameters))
		}
		for key, actual := range testMap {
			expected := testMap[key]
			if len(expected) != len(actual) {
				t.Fatalf("error adding query parameters from struct: different number of expected and actual (%d != %d)", len(expected), len(actual))
			}
			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Fatalf("error adding query parameters from struct: different values for %s: expected %s, got %s", key, expected[i], actual[i])
				}
			}
		}
	}

	defer handler("only structs and maps can be passed as sources", t)
	New("").Add().QueryParametersFrom(&s)
}

func TestAddHeader(t *testing.T) {
	f := New("").Add().Header("key", "value1", "value2")
	if len(f.headers) != 1 {
		t.Fatalf("error adding headers: expected 1, got %d", len(f.headers))
	}
	if len(f.headers["Key"]) != 2 {
		t.Fatalf("error adding headers: expected 2, got %d", len(f.headers["Key"]))
	}
	if f.headers["Key"][0] != "value1" {
		t.Fatalf("error adding headers: expected \"value1\", got %q", f.headers["Key"][0])
	}
	if f.headers["Key"][1] != "value2" {
		t.Fatalf("error adding headers: expected \"value2\", got %q", f.headers["Key"][1])
	}
	// now add one more
	f.Header("key", "value3")
	if len(f.headers) != 1 {
		t.Fatalf("error adding headers: expected 1, got %d", len(f.headers))
	}
	if len(f.headers["Key"]) != 3 {
		t.Fatalf("error adding headers: expected 3, got %d", len(f.headers["Key"]))
	}
	if f.headers["Key"][0] != "value1" {
		t.Fatalf("error adding headers: expected \"value1\", got %q", f.headers["Key"][0])
	}
	if f.headers["Key"][1] != "value2" {
		t.Fatalf("error adding headers: expected \"value2\", got %q", f.headers["Key"][1])
	}
	if f.headers["Key"][2] != "value3" {
		t.Fatalf("error adding headers: expected \"value3\", got %q", f.headers["Key"][2])
	}
}

func TestSetHeader(t *testing.T) {
	f := New("").Add().Header("key", "value1", "value2")
	f.Set().Header("key", "value0", "value1", "value2", "value3")
	if len(f.headers) != 1 {
		t.Fatalf("error setting headers: expected 1, got %d", len(f.headers))
	}
	if len(f.headers["Key"]) != 4 {
		t.Fatalf("error setting headers: expected 4, got %d", len(f.headers["Key"]))
	}
	if f.headers["Key"][0] != "value0" {
		t.Fatalf("error setting headers: expected \"value0\", got %q", f.headers["Key"][0])
	}
	if f.headers["Key"][1] != "value1" {
		t.Fatalf("error setting headers: expected \"value1\", got %q", f.headers["Key"][1])
	}
	if f.headers["Key"][2] != "value2" {
		t.Fatalf("error setting headers: expected \"value2\", got %q", f.headers["Key"][2])
	}
	if f.headers["Key"][3] != "value3" {
		t.Fatalf("error setting headers: expected \"value3\", got %q", f.headers["Key"][3])
	}
}

func TestDelHeader(t *testing.T) {
	f := New("").Add().Header("key", "value1", "value2")
	f.Del().Header("Key")
	if len(f.headers) != 0 {
		t.Fatalf("error deleting headers: expected 0, got %d", len(f.headers))
	}
}

func TestRemHeader(t *testing.T) {
	f := New("").
		Add().
		Header("key1", "value1", "value2").
		Header("key2", "value1", "value2").
		Header("key3", "value1", "value2").
		Header("another_key", "value1", "value2")
	f.Remove().Header("^Key\\d$")
	if len(f.headers) != 1 {
		t.Fatalf("error removing multiple headers: expected 1, got %d", len(f.headers))
	}
}

func TestHeadersFrom(t *testing.T) {

	type Nested struct {
		Query7 string `header:"header7"`
		Query1 string `header:"header1"`
	}

	type Embedded struct {
		Query5 string `header:"header5"`
		Query6 string `header:"header6"`
	}

	type Struct struct {
		Query1 string  `header:"header1"`
		Query2 *string `header:"header2"`
		Query3 bool    `header:"header3"`
		Query4 *bool   `header:"header4"`
		Embedded
		Nested *Nested
	}

	s := "value2"
	b := true
	test := &Struct{
		Query1: "value1a",
		Query2: &s,
		Query3: true,
		Query4: &b,
		Embedded: Embedded{
			Query5: "value5",
			Query6: "value6",
		},
		Nested: &Nested{
			Query7: "value7",
			Query1: "value1b",
		},
	}

	f := New("").Add().HeadersFrom(test)
	if len(f.headers) != 7 {
		t.Fatalf("error adding headers from struct: expected 7, got %d", len(f.headers))
	}

	tests := map[string][]string{
		"Header1": []string{"value1a", "value1b"},
		"Header2": []string{"value2"},
		"Header3": []string{"true"},
		"Header4": []string{"true"},
		"Header5": []string{"value5"},
		"Header6": []string{"value6"},
		"Header7": []string{"value7"},
	}
	for key, actual := range tests {
		expected := tests[key]
		if len(expected) != len(actual) {
			t.Fatalf("error adding headers from struct: different number of expected and actual (%d != %d)", len(expected), len(actual))
		}
		for i := 0; i < len(expected); i++ {
			if expected[i] != actual[i] {
				t.Fatalf("error adding headers from struct: different values for %s: expected %s, got %s", key, expected[i], actual[i])
			}
		}
	}
}

func TestWithEntity(t *testing.T) {
	expected := "some text to send along"
	f := New("").ContentType("text/plain").WithEntity(strings.NewReader(expected))
	data, _ := ioutil.ReadAll(f.body)
	actual := string(data)
	if actual != expected {
		t.Fatalf("error adding entity by reader: expected %s, got %s", expected, actual)
	}
}

func TestWithJSONEntity(t *testing.T) {
	type A struct {
		Field1 string  `json:"field1,omitempty"`
		Field2 bool    `json:"field2,omitempty"`
		Field3 int     `json:"field3,omitempty"`
		Field4 *string `json:"field4,omitempty"`
	}

	s := "value4"
	a := A{
		Field1: "value1",
		Field2: true,
		Field3: 12,
		Field4: &s,
	}
	expected := "{\"field1\":\"value1\",\"field2\":true,\"field3\":12,\"field4\":\"value4\"}"

	// test with struct "by value"
	f := New("").WithJSONEntity(a)
	data, _ := ioutil.ReadAll(f.body)
	actual := string(data)
	if actual != expected {
		t.Fatalf("error adding entity by reader: expected %s, got %s", expected, actual)
	}
	if f.headers["Content-Type"][0] != "application/json" {
		t.Fatalf("error adding entity by reader: content type is %s, expected \"application/json\"", f.headers["Content-Type"][0])
	}

	f = New("").ContentType("application/my-type").WithJSONEntity(&a)
	data, _ = ioutil.ReadAll(f.body)
	actual = string(data)
	if actual != expected {
		t.Fatalf("error adding entity by reader: expected %s, got %s", expected, actual)
	}
	if f.headers["Content-Type"][0] != "application/my-type" {
		t.Fatalf("error adding entity by reader: content type is %s, expected \"application/my-type\"", f.headers["Content-Type"][0])
	}
}

func TestWithJSONEntityNoStruct(t *testing.T) {
	defer handler("only structs can be passed as source for JSON entities", t)
	s := "a string"
	New("").WithJSONEntity(s)
}

func TestWithJSONEntityNoStructPtr(t *testing.T) {
	defer handler("only structs can be passed as source for JSON entities", t)
	s := "a string"
	New("").WithJSONEntity(&s)
}

func TestWithXMLEntity(t *testing.T) {
	type A struct {
		Field1 string  `xml:"field1,omitempty"`
		Field2 bool    `xml:"field2,omitempty"`
		Field3 int     `xml:"field3,omitempty"`
		Field4 *string `xml:"field4,omitempty"`
	}

	s := "value4"
	a := A{
		Field1: "value1",
		Field2: true,
		Field3: 12,
		Field4: &s,
	}
	expected := "<A><field1>value1</field1><field2>true</field2><field3>12</field3><field4>value4</field4></A>"

	// test with struct "by value"
	f := New("").WithXMLEntity(a)
	data, _ := ioutil.ReadAll(f.body)
	actual := string(data)
	if actual != expected {
		t.Fatalf("error adding entity by reader: expected %s, got %s", expected, actual)
	}
	if f.headers["Content-Type"][0] != "text/xml" {
		t.Fatalf("error adding entity by reader: content type is %s, expected \"text/xml\"", f.headers["Content-Type"][0])
	}

	f = New("").ContentType("application/my-type").WithXMLEntity(&a)
	data, _ = ioutil.ReadAll(f.body)
	actual = string(data)
	if actual != expected {
		t.Fatalf("error adding entity by reader: expected %s, got %s", expected, actual)
	}
	if f.headers["Content-Type"][0] != "application/my-type" {
		t.Fatalf("error adding entity by reader: content type is %s, expected \"application/my-type\"", f.headers["Content-Type"][0])
	}
}

func TestWithXMLEntityNoStruct(t *testing.T) {
	defer handler("only structs can be passed as source for XML entities", t)
	s := "a string"
	New("").WithXMLEntity(s)
}

func TestWithXMLEntityNoStructPtr(t *testing.T) {
	defer handler("only structs can be passed as source for XML entities", t)
	s := "a string"
	New("").WithXMLEntity(&s)
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
