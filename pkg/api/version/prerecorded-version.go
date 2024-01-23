// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package handles the versioning in the API for prerecorded endpoint
package version

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/go-querystring/query"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/common"
)

const (
	// PrerecordedAPIVersion current supported version
	PrerecordedAPIVersion string = "v1"

	// PrerecordedPath is the current path for prerecorded transcription
	PrerecordedPath string = "listen"
)

/*
GetPrerecordedAPI is a function which controls the versioning of the live transcription API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for the live transcription
*/
func GetPrerecordedAPI(ctx context.Context, host, version, path string, options interfaces.PreRecordedTranscriptionOptions, args ...interface{}) (string, error) {
	if path == "" {
		path = PrerecordedPath
	}

	// handle protocol and host
	protocol := APIProtocol
	if host == "" {
		host = common.DefaultHost
	}

	r, err := regexp.Compile("^(wss|ws)://(.+)$")
	if err != nil {
		klog.V(1).Infof("regexp.Compile err: %v\n", err)
		return "", err
	}

	match := r.MatchString(host)
	klog.V(3).Infof("host decompose... match: %t\n", match)
	if match {
		matches := r.FindStringSubmatch(host)
		for _, match := range matches {
			klog.V(5).Infof("match: %s\n", match)
		}
		protocol = matches[1]
		host = matches[2]

		slash := strings.Index(host, "/")
		if slash > 0 {
			host = host[:slash]
		}
	}
	klog.V(3).Infof("protocol: %s\n", protocol)
	klog.V(3).Infof("host: %s\n", host)

	// handle version and path
	if version == "" {
		version = PrerecordedAPIVersion
	}

	r, err = regexp.Compile("^(v[0-9]+|%%s)/")
	if err != nil {
		klog.V(1).Infof("regexp.Compile err: %v\n", err)
		return "", err
	}

	match = r.MatchString(path)
	klog.V(3).Infof("path decompose - match: %t\n", match)
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
	u := url.URL{Scheme: protocol, Host: host, Path: completeFullpath, RawQuery: q.Encode()}

	return u.String(), nil
}
