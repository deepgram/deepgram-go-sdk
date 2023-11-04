// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

// LiveMessageCallback is a callback used to receive notifcations for platforms messages
type LiveMessageCallback interface {
	Message(mr *MessageResponse) error
	Metadata(md *MetadataResponse) error
	// TODO: implement other conversation insights
}
