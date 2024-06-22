// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

// These are the message types that can be received from the text-to-speech streaming API
const (
	// message types
	TypeOpenResponse     string = "Open"
	TypeMetadataResponse string = "Metadata"
	TypeFlushedResponse  string = "Flushed"
	TypeCloseResponse    string = "Close"

	// Error type
	TypeErrorResponse string = "Error"
)
