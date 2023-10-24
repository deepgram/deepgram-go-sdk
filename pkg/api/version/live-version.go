// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package handles the versioning in the API both async and streaming
*/
package version

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

const (
	// version
	LiveAPIVersion string = "v1"

	// paths
	LivePath string = "%s/listen"
)

func GetLiveAPI(ctx context.Context, options interfaces.LiveTranscriptionOptions) (string, error) {
	if options.Host == "" {
		options.Host = DefaultHost
	}
	if options.ApiVersion == "" {
		options.ApiVersion = LiveAPIVersion
	}

	q, err := query.Values(options)
	if err != nil {
		return "", err
	}

	if parameters, ok := ctx.Value(interfaces.ParametersContext{}).(map[string][]string); ok {
		for k, vs := range parameters {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
	}

	u := url.URL{Scheme: "wss", Host: options.Host, Path: fmt.Sprintf(LivePath, options.ApiVersion), RawQuery: q.Encode()}
	return u.String(), nil
}
