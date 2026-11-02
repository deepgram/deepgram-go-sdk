// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"errors"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
)

func Test_TransportErrorsSetCode(t *testing.T) {
	for _, test := range []struct {
		name string
		err  error
		code string
	}{
		{name: "close reason", err: errors.New("websocket: close 1008 (policy violation): DATA-0000 invalid configuration"), code: "DATA-0000"},
		{name: "bad handshake", err: errors.New("websocket: bad handshake")},
	} {
		t.Run(test.name, func(t *testing.T) {
			for name, response := range map[string]*msginterfaces.ErrorResponse{
				"callback": (&WSCallback{}).errorToResponse(test.err),
				"channel":  (&WSChannel{}).errorToResponse(test.err),
			} {
				t.Run(name, func(t *testing.T) {
					if response.Code != test.code {
						t.Errorf("Code = %q, want %q", response.Code, test.code)
					}
				})
			}
		})
	}
}
