// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package restv1

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"sync/atomic"
	"time"

	"github.com/deepgram/deepgram-go-sdk/v3/pkg/client/rest/v1/debug"
)

var (
	// Trace reads an http request or response from rc and writes to w.
	// The content type (kind) should be one of "xml" or "json".
	Trace = func(rc io.ReadCloser, w io.Writer, kind string) io.ReadCloser {
		return debug.NewTeeReader(rc, w)
	}
)

// debugRoundTrip contains state and logic needed to debug a single round trip.
type debugRoundTrip struct {
	cn  uint64         // Client number
	rn  uint64         // Request number
	log io.WriteCloser // Request log
	cs  []io.Closer    // Files that need closing when done
}

func (d *debugRoundTrip) logf(format string, a ...interface{}) {
	now := time.Now().Format("2006-01-02T15-04-05.000000000")
	fmt.Fprintf(d.log, "%s - %04d: ", now, d.rn)
	fmt.Fprintf(d.log, format, a...)
	fmt.Fprintf(d.log, "\n")
}

func (d *debugRoundTrip) enabled() bool {
	return d != nil
}

func (d *debugRoundTrip) done() {
	for _, c := range d.cs {
		c.Close()
	}
}

func (d *debugRoundTrip) newFile(suffix string) io.WriteCloser {
	return debug.NewFile(fmt.Sprintf("%d-%04d.%s", d.cn, d.rn, suffix))
}

func (d *debugRoundTrip) ext(h http.Header) string {
	const json = "application/json"
	ext := "xml"
	if h.Get("Accept") == json || h.Get("Content-Type") == json {
		ext = "json"
	}
	return ext
}

func (d *debugRoundTrip) debugRequest(req *http.Request) string {
	if d == nil {
		return ""
	}

	// Capture headers
	wc := d.newFile("req.headers")
	b, _ := httputil.DumpRequest(req, false)
	_, err := wc.Write(b)
	if err != nil {
		d.logf("Error writing request headers: %v", err)
	}
	wc.Close()

	ext := d.ext(req.Header)
	// Capture body
	wc = d.newFile("req." + ext)
	if req.Body != nil {
		req.Body = Trace(req.Body, wc, ext)
	}

	// Delay closing until marked done
	d.cs = append(d.cs, wc)

	return ext
}

func (d *debugRoundTrip) debugResponse(res *http.Response, ext string) {
	if d == nil {
		return
	}

	// Capture headers
	wc := d.newFile("res.headers")
	b, _ := httputil.DumpResponse(res, false)
	_, err := wc.Write(b)
	if err != nil {
		d.logf("Error writing response headers: %v", err)
	}
	wc.Close()

	// Capture body
	wc = d.newFile("res." + ext)
	res.Body = Trace(res.Body, wc, ext)

	// Delay closing until marked done
	d.cs = append(d.cs, wc)
}

var cn uint64 // Client counter

// debugContainer wraps the debugging state for a single client.
type debugContainer struct {
	cn  uint64         // Client number
	rn  uint64         // Request counter
	log io.WriteCloser // Request log
}

func newDebug() *debugContainer {
	d := debugContainer{
		cn: atomic.AddUint64(&cn, 1),
		rn: 0,
	}

	if v := os.Getenv("DEEPGRAM_DEBUG_REST"); v != "" {
		fmt.Printf("\n---------------------------------\n")
		fmt.Printf("DEEPGRAM_DEBUG_REST found\n")
		fmt.Printf("---------------------------------\n\n")
		debug.SetProvider(debug.LogProvider{})
	}

	if !debug.Enabled() {
		return nil
	}

	d.log = debug.NewFile(fmt.Sprintf("%d-client.log", d.cn))
	return &d
}

func (d *debugContainer) newRoundTrip() *debugRoundTrip {
	if d == nil {
		return nil
	}

	drt := debugRoundTrip{
		cn:  d.cn,
		rn:  atomic.AddUint64(&d.rn, 1),
		log: d.log,
	}

	return &drt
}
