// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	"net/http"
	"net/url"
	"testing"
)

type compatibleClientOptions struct{}

func (*compatibleClientOptions) Parse() error                                    { return nil }
func (*compatibleClientOptions) SetAPIKey(string)                                {}
func (*compatibleClientOptions) GetAuthToken() (string, bool)                    { return "", false }
func (*compatibleClientOptions) GetHost() string                                 { return "" }
func (*compatibleClientOptions) GetWSHeaderProcessor() func(http.Header)         { return nil }
func (*compatibleClientOptions) GetRedirectService() bool                        { return false }
func (*compatibleClientOptions) GetSkipServerAuth() bool                         { return false }
func (*compatibleClientOptions) GetProxy() func(*http.Request) (*url.URL, error) { return nil }

func Test_ClientOptionsDoesNotRequireReconnectPolicy(t *testing.T) {
	var options ClientOptions = &compatibleClientOptions{}
	if options == nil {
		t.Fatal("legacy client options did not satisfy ClientOptions")
	}
}
