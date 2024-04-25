// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

// MessageType is the header to bootstrap you way unmarshalling other messages
/*
	Example:
	{
		"type": "message",
		"message": {
			...
		}
	}
*/
type MessageType struct {
	Type string `json:"type"`
}
