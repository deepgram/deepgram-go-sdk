// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import "testing"

const emptyAlternativesResults = `{"type":"Results","channel":{"alternatives":[]},"is_final":false}`

func Test_InspectEmptyAlternatives(t *testing.T) {
	t.Run("callback", func(t *testing.T) {
		client := &WSCallback{}
		if err := client.inspect([]byte(emptyAlternativesResults)); err != nil {
			t.Fatalf("inspect returned an error: %v", err)
		}
		if client.lastDatagram != nil {
			t.Fatal("lastDatagram was set for a message without a transcript")
		}
	})
	t.Run("channel", func(t *testing.T) {
		client := &WSChannel{}
		if err := client.inspect([]byte(emptyAlternativesResults)); err != nil {
			t.Fatalf("inspect returned an error: %v", err)
		}
		if client.lastDatagram != nil {
			t.Fatal("lastDatagram was set for a message without a transcript")
		}
	})
}

func Test_InspectInterimSetsLastDatagram(t *testing.T) {
	const interim = `{"type":"Results","channel":{"alternatives":[{"transcript":"hello"}]},"is_final":false}`

	callback := &WSCallback{}
	if err := callback.inspect([]byte(interim)); err != nil {
		t.Fatalf("callback inspect returned an error: %v", err)
	}
	if callback.lastDatagram == nil {
		t.Fatal("callback lastDatagram was not set for an interim transcript")
	}

	channel := &WSChannel{}
	if err := channel.inspect([]byte(interim)); err != nil {
		t.Fatalf("channel inspect returned an error: %v", err)
	}
	if channel.lastDatagram == nil {
		t.Fatal("channel lastDatagram was not set for an interim transcript")
	}
}
