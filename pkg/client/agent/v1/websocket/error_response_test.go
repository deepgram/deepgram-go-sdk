// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"errors"
	"testing"
)

func Test_TransportErrorsSetCode(t *testing.T) {
	response := (&WSChannel{}).errorToResponse(errors.New("websocket: bad 400 handshake: invalid configuration"))
	if response.Code != "bad" {
		t.Errorf("Code = %q, want %q", response.Code, "bad")
	}
}
