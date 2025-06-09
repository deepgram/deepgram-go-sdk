// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/schema"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/v3/pkg/common"
)

// getAPIURL constructs the URL for API requests and handles versioning, path resolution,
// host overrides, and query string parameters.
func getAPIURL(ctx context.Context, apiType, host, version, path string, options interface{}, args ...interface{}) (string, error) {
	// set the path
	if path == "" {
		path = APIPathMap[apiType]
	}

	// set the protocol
	protocol := HTTPProtocol
	if apiType == APITypeLive || apiType == APITypeSpeakStream || apiType == APITypeAgent {
		protocol = WSProtocol
	}
	if host == "" {
		host = common.DefaultHost
	}
	if apiType == APITypeAgent && host == common.DefaultHost {
		klog.V(4).Infof("overriding with agent host\n")
		host = common.DefaultAgentHost
	}

	// check if the host has a protocol
	r := regexp.MustCompile(`^(https?)://(.+)$`)
	if apiType == APITypeLive || apiType == APITypeSpeakStream {
		r = regexp.MustCompile(`^(wss?)://(.+)$`)
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

		if slash := strings.Index(host, "/"); slash > 0 {
			host = host[:slash]
		}
	}
	klog.V(3).Infof("protocol: %s\n", protocol)
	klog.V(3).Infof("host: %s\n", host)

	// make sure the version is set
	if version == "" {
		version = DefaultAPIVersion
	}

	// remove the version from the path if it exists
	r = regexp.MustCompile(`^(v\d+|%%s)/`)
	match = r.MatchString(path)
	klog.V(3).Infof("path decompose - match: %t\n", match)
	if match {
		path = r.ReplaceAllString(path, "")
	}

	// Construct the query string
	// using an options struct
	var q url.Values
	if options != nil {
		q = url.Values{}

		encoder := schema.NewEncoder()
		if err := encoder.Encode(options, q); err != nil {
			return "", err
		}
	}

	// using custom parameters
	if parameters, ok := ctx.Value(interfaces.ParametersContext{}).(map[string][]string); ok {
		for k, vs := range parameters {
			for _, v := range vs {
				if q == nil {
					q = url.Values{}
				}
				q.Add(k, v)
			}
		}
	}

	// construct the full path and substitute the version and all query parameters
	fullpath := fmt.Sprintf("%%s/%s", path)
	completeFullpath := fmt.Sprintf(fullpath, append([]interface{}{version}, args...)...)
	klog.V(3).Infof("completeFullpath: %s\n", completeFullpath)

	// construct the URL
	var u url.URL
	if q != nil {
		u = url.URL{Scheme: protocol, Host: host, Path: completeFullpath, RawQuery: q.Encode()}
	} else {
		u = url.URL{Scheme: protocol, Host: host, Path: completeFullpath}
	}
	klog.V(3).Infof("URI final: %s\n", u.String())

	return u.String(), nil
}
