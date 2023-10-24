// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package prerecorded

import (
	rest "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/rest"
)

// Client implements helper functionality for Prerecorded API
type Client struct {
	*rest.Client
}
