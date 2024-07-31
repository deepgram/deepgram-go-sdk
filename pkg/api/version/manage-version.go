// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package version

import (
	"context"
)

const (
	// version
	ManageAPIVersion string = "v1"

	// balances
	BalancesURI     string = "projects/%s/balances"
	BalancesByIDURI string = "projects/%s/balances/%s"

	// invitations
	InvitationsURI      string = "projects/%s/invites"
	InvitationsByIDURI  string = "projects/%s/invites/%s"
	InvitationsLeaveURI string = "projects/%s/leave"

	// Keys
	KeysURI     string = "projects/%s/keys"
	KeysByIDURI string = "projects/%s/keys/%s"

	// Members
	MembersURI     string = "projects/%s/members"
	MembersByIDURI string = "projects/%s/members/%s"

	// models
	ModelsURI            string = "models"
	ModelsByIDURI        string = "models/%s"
	ModelsProjectURI     string = "projects/%s/models"
	ModelsProjectByIDURI string = "projects/%s/models/%s"

	// projects
	ProjectsURI     string = "projects"
	ProjectsByIDURI string = "projects/%s"

	// scopes
	MembersScopeByIDURI string = "projects/%s/members/%s/scopes"

	// usage
	UsageRequestURI     string = "projects/%s/requests"
	UsageRequestByIDURI string = "projects/%s/requests/%s"
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
	return getAPIURL(ctx, "speak", host, version, path, nil, args...)
}
