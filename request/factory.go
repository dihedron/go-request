// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"io"
	"net/http"
)

// Factory is the HTTP request factory.
type Factory struct {
	// client is the HTTP client that will make the requests.
	client http.Client
	// method is the request's HTTP method (GET, POST, etc.).
	method string
	// raw url string for requests
	url string

	headers HeadersProvider

	parameters ParametersProvider

	entity EntityProvider
}

type Parameters map[string]interface{}

type ParametersProvider interface {
	ProvideParameters() Parameters
}

type Entity struct {
	ContentType string
	Reader      io.Reader
}

type EntityProvider interface {
	ProvideEntity() Entity
}
