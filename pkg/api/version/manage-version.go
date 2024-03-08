// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	common "github.com/deepgram/deepgram-go-sdk/pkg/common"
	"github.com/google/go-querystring/query"
)

const (
	// version
	ManageAPIVersion string = "v1"

	// balances
	BalancesURI     string = "projects/%s/balances"
	BalancesByIdURI string = "projects/%s/balances/%s"

	// invitations
	InvitationsURI      string = "projects/%s/invites"
	InvitationsByIdURI  string = "projects/%s/invites/%s"
	InvitationsLeaveURI string = "projects/%s/leave"

	// Keys
	KeysURI     string = "projects/%s/keys"
	KeysByIdURI string = "projects/%s/keys/%s"

	// Members
	MembersURI     string = "projects/%s/members"
	MembersByIdURI string = "projects/%s/members/%s"

	// projects
	ProjectsURI     string = "projects"
	ProjectsByIdURI string = "projects/%s"

	// scopes
	MembersScopeByIdURI string = "projects/%s/members/%s/scopes"

	// usage
	UsageRequestURI     string = "projects/%s/requests"
	UsageRequestByIdURI string = "projects/%s/requests/%s"
	UsageURI            string = "projects/%s/usage"
	UsageFieldsURI      string = "projects/%s/usage/fields"
)

/*
GetManageAPI is a function which controls the versioning of the manage API and provides
mechanism for:

- overriding the host endpoint
- overriding the version used
- overriding the endpoint path
- additional arguments to the query string/parameters

The return value is the complete URL endpoint to be used for manage
*/
func GetManageAPI(ctx context.Context, host, version, path string, vals interface{}, args ...interface{}) (string, error) {
	if path == "" {
		return "", ErrInvalidPath
	}

	// handle protocol and host
	protocol := APIProtocol
	if host == "" {
		host = common.DefaultHost
	}

	r, err := regexp.Compile("^(https|http)://(.+)$")
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
		version = ManageAPIVersion
	}

	r, err = regexp.Compile("^(v[0-9]+|%%s)/")
	if err != nil {
		klog.V(1).Infof("regexp.Compile err: %v\n", err)
		return "", err
	}

	match = r.MatchString(path)
	klog.V(3).Infof("match: %t\n", match)

	if match {
		// version = r.FindStringSubmatch(path)[0]
		path = r.ReplaceAllString(path, "")
	}

	var q url.Values
	if vals != nil {
		q, err = query.Values(vals)
		if err != nil {
			return "", err
		}
	}

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

	fullpath := fmt.Sprintf("%%s/%s", path)
	completeFullpath := fmt.Sprintf(fullpath, append([]interface{}{version}, args...)...)

	var u url.URL
	if q != nil {
		u = url.URL{Scheme: protocol, Host: host, Path: completeFullpath, RawQuery: q.Encode()}
	} else {
		u = url.URL{Scheme: protocol, Host: host, Path: completeFullpath}
	}

	return u.String(), nil
}
