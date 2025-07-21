// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

/*
This package provides the types for the Deepgram Auth API.
*/
package interfaces

/***********************************/
// shared/common structs
/***********************************/

// GrantTokenRequest represents the request body for the grant token endpoint
type GrantTokenRequest struct {
	TTLSeconds *int `json:"ttl_seconds,omitempty"`
}

// GrantToken provides a JWT
type GrantToken struct {
	AccessToken string  `json:"access_token,omitempty"`
	ExpiresIn   float64 `json:"expires_in,omitempty"`
}
