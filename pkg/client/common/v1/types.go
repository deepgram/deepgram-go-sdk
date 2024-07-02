// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv1

import (
	rest "github.com/deepgram/deepgram-go-sdk/pkg/client/rest/v1"
)

type RESTClient = rest.Client

// Client implements helper functionality for Prerecorded API
type Client struct {
	*RESTClient
}
