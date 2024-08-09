// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

// SpeakMessageCallback is a callback used to receive notifications for platforms messages
type SpeakMessageCallback interface {
	// These are WS TextMessage that are used for flow control.
	Open(or *OpenResponse) error
	Metadata(md *MetadataResponse) error
	Flush(fl *FlushedResponse) error
	Close(cr *CloseResponse) error

	Warning(er *WarningResponse) error
	Error(er *ErrorResponse) error

	// These are WS BinaryMessage that are used to send audio data to the client
	Binary(byMsg []byte) error
}
