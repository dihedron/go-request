// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import "io"

// ValueProvider is a provider capable or returning a map of keys to values;
// these can be used for setting query parameters and headers in an HTTP request.
type ValueProvider interface {
	// Provide returns a string-to-string-array map; the keys and values in the
	// returned map will be set as request headers or query parameters.
	Provide() map[string][]string
}

// Entity represents the HTTP entity in a request; it will be used to populate
// the request body with data.
type Entity struct {
	ContentType string
	Reader      io.Reader
}

// EntityProvider is the common interface to structs providing HTTP request
// payload.
type EntityProvider interface {
	// Provide returns an HTTP request ody generator: the returned object can be
	// used to retrieve the content type and the data associated with the request
	// entity. The returned Entity can be nil.
	Provide() *Entity
}
