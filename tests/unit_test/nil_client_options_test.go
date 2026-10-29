// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"testing"

	analyze "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/analyze"
	auth "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth"
	listen "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
	manage "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/manage"
	speak "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak"
)

func Test_RESTConstructorsAllowNilClientOptions(t *testing.T) {
	constructors := map[string]func() interface{}{
		"analyze": func() interface{} { return analyze.New("test-api-key", nil) },
		"auth":    func() interface{} { return auth.New("test-api-key", nil) },
		"listen":  func() interface{} { return listen.NewREST("test-api-key", nil) },
		"manage":  func() interface{} { return manage.New("test-api-key", nil) },
		"speak":   func() interface{} { return speak.NewREST("test-api-key", nil) },
	}

	for name, newClient := range constructors {
		t.Run(name, func(t *testing.T) {
			if client := newClient(); client == nil {
				t.Fatal("New returned nil for nil ClientOptions")
			}
		})
	}
}
