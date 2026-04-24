// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv2

import (
	commonv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
)

const (
	PackageVersion string = "v2.0"
)

// WSClient is the v2 WebSocket client. It embeds v1's WSClient and adds WritePing.
type WSClient struct {
	*commonv1.WSClient
}
