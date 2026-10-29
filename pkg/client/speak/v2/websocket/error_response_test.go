// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv2

import (
	"errors"
	"testing"

	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
)

func Test_TransportErrorsSetCode(t *testing.T) {
	err := errors.New("websocket: bad 400 handshake: invalid configuration")
	for name, response := range map[string]*msginterfaces.ErrorResponse{
		"callback": (&WSCallback{}).errorToResponse(err),
		"channel":  (&WSChannel{}).errorToResponse(err),
	} {
		t.Run(name, func(t *testing.T) {
			if response.Code != "bad" {
				t.Errorf("Code = %q, want %q", response.Code, "bad")
			}
		})
	}
}
