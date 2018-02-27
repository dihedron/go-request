// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"net/http"
)

// Factory is the HTTP request factory.
type Factory struct {
	// client is the HTTP client that will make the requests.
	Client *http.Client

	// baseURL is the base UTL for generating relative URL requests.
	BaseURL string

	Headers []*ValueProvider

	Parameters []*ValueProvider

	Entity *EntityProvider
}

func NewFactory(client *http.Client, baseURL string) *Factory {
	return &Factory{
		Client:     client,
		BaseURL:    baseURL,
		Headers:    []*ValueProvider{},
		Parameters: []*ValueProvider{},
	}
}
