// Copyright 2017-present Andrea Funtò. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package request

import (
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/dihedron/go-log/log"
)

func TestFactory(t *testing.T) {
	var client = &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}
	request, _ := New("GET", "http://www.repubblica.it").UserAgent("MyCrawler/1.0").Make()
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("no network connection: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("error reading response body: %v", err)
	}
	log.Debugf("response:\n%s", body)
}
