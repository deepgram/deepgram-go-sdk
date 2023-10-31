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
	"regexp"

	"github.com/google/go-querystring/query"

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	common "github.com/deepgram-devs/deepgram-go-sdk/pkg/common"
)

const (
	// version
	LiveAPIVersion string = "v1"

	// paths
	LivePath string = "listen"
)

func GetLiveAPI(ctx context.Context, host, version, path string, options interfaces.LiveTranscriptionOptions, args ...interface{}) (string, error) {
	if path == "" {
		return "", ErrInvalidPath
	}

	if host == "" {
		host = common.DefaultHost
	}
	if version == "" {
		version = LiveAPIVersion
	}

	r, err := regexp.Compile("^(v[0-9]+|%%s)/")
	if err != nil {
		// fmt.Printf("regexp.Compile err: %v\n", err)
		return "", err
	}

	match := r.MatchString(path)
	fmt.Printf("match: %t\n", match)

	if match {
		// version = r.FindStringSubmatch(path)[0]
		path = r.ReplaceAllString(path, "")
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

	fullpath := fmt.Sprintf("%%s/%s", path)
	completeFullpath := fmt.Sprintf(fullpath, append([]interface{}{version}, args...)...)
	u := url.URL{Scheme: "wss", Host: host, Path: completeFullpath, RawQuery: q.Encode()}

	return u.String(), nil
}
