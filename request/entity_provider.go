// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

// JSONProvider is a specialisation of EntityProvider for JSON payloads.
type JSONProvider struct {
	Entity interface{}
}

func (p *JSONProvider) Provide() *Entity {
	return nil
}
