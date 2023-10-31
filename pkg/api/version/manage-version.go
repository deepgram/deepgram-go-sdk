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
)

const (
	// version
	ManageAPIVersion string = "v1"

	// Project: /v1/projects
	// Keys: /v1/projects/<project_id>/keys
	// Members: /v1/projects/members

	// keys
	// Path:              "/v1/projects",
	// TranscriptionPath: "/v1/listen",
	//
	// path := fmt.Sprintf("%s/%s/keys", dg.Client.Path, projectId) //list
	// path := fmt.Sprintf("%s/%s/keys/%s", dg.Client.Path, projectId, keyId) //get
	// path := fmt.Sprintf("%s/%s/keys", dg.Client.Path, projectId) // create
	// path := fmt.Sprintf("%s/%s/keys/%s", dg.Client.Path, projectId, keyId) //delete

)

func GetManageAPI(ctx context.Context, host, version, path, key string, vals interface{}) (string, error) {
	if path == "" {
		return "", ErrInvalidPath
	}

	if host == "" {
		host = DefaultHost
	}
	if version == "" {
		version = ManageAPIVersion
	}

	r, err := regexp.Compile("^(v[0-9]+)/")
	if err != nil {
		fmt.Printf("regexp.Compile err: %v\n", err)
		return "", err
	}

	match := r.MatchString(path)
	fmt.Printf("match: %t\n", match)

	if match {
		// version = r.FindStringSubmatch(path)[0]
		path = r.ReplaceAllString(path, "")
	}

	q, err := query.Values(vals)
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

	u := url.URL{Scheme: "https", Host: host, Path: fmt.Sprintf("%s/%s", version, path), RawQuery: q.Encode()}
	return u.String(), nil
}
