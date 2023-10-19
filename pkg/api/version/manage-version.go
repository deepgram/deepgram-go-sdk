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

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	common "github.com/deepgram-devs/deepgram-go-sdk/pkg/common"
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

func GetManageAPI(ctx context.Context, host, version, path string, vals interface{}, args ...interface{}) (string, error) {
	if path == "" {
		return "", ErrInvalidPath
	}

	if host == "" {
		host = common.DefaultHost
	}
	if version == "" {
		version = ManageAPIVersion
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
		u = url.URL{Scheme: "https", Host: host, Path: completeFullpath, RawQuery: q.Encode()}
	} else {
		u = url.URL{Scheme: "https", Host: host, Path: completeFullpath}
	}

	return u.String(), nil
}
