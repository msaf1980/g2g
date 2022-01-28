// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package expvar provides a standardized interface to public variables, such
// as operation counters in servers. It exposes these variables via HTTP at
// /debug/vars in JSON format.
//
// Operations to set or modify these public variables are atomic.
//
// In addition to adding the HTTP handler, this package registers the
// following variables:
//
//	cmdline   os.Args
//	memstats  runtime.Memstats
//
// The package is sometimes only imported for the side effect of
// registering its HTTP handler and the above variables. To use it
// this way, link this package into your program:
//	import _ "expvar"
//
package g2g

import (
	"strconv"
	"testing"
)

func TestRoundFloat(t *testing.T) {
	m := map[float64]string{
		0.00:  "0",
		123.0: "123",
		1.2:   "1.2",

		1.00:        "1",
		1.001:       "1.001",
		1.00000001:  "1.00000001",
		0.00001:     "0.00001",
		0.01000:     "0.01",
		0.01999:     "0.01999",
		-1.234:      "-1.234",
		123.456:     "123.456",
		99999.09123: "99999.09123",
	}
	for v, expected := range m {
		if got := RoundFloat(v); got != expected {
			t.Errorf("%s: got %s, expected %s", strconv.FormatFloat(v, 'g', -1, 64), got, expected)
		}
	}
}
