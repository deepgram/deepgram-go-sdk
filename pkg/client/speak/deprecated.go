// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package speak

import (
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	speakv1rest "github.com/deepgram/deepgram-go-sdk/pkg/client/speak/v1/rest"
)

/***********************************/
// Deprecated (THESE WILL STILL WORK,
// BUT WILL BE REMOVED IN A FUTURE RELEASE)
/***********************************/
/*
Legacy Client Name

Deprecated: This struct is deprecated. Please use RestClient struct. This will be removed in a future release.
*/
type Client = speakv1rest.RESTClient

/*
NewWithDefaults creates a new speak client with all default options

Deprecated: This function is deprecated. Please use NewREST(). This will be removed in a future release.
*/
func NewWithDefaults() *speakv1rest.RESTClient {
	return speakv1rest.NewWithDefaults()
}

/*
New creates a new speak client with the specified options

Input parameters:
- ctx: context.Context object
- apiKey: string containing the Deepgram API key
- options: ClientOptions which allows overriding things like hostname, version of the API, etc.

Deprecated: This function is deprecated. Please use NewREST(). This will be removed in a future release.
*/
func New(apiKey string, options *interfaces.ClientOptions) *speakv1rest.RESTClient {
	return speakv1rest.New(apiKey, options)
}
