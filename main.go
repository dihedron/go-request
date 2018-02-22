// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/dihedron/go-log/log"
)

func main() {

	log.SetStream(os.Stdout, true)
	log.Debugln("hallo")
}
