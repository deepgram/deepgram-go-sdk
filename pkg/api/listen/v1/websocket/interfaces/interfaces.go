// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines interfaces for the live API
package interfacesv1

/*
Chan Interfaces
*/
// LiveMessageCallback is a callback used to receive notifcations for platforms messages
type LiveMessageChan interface {
	GetOpen() []*chan *OpenResponse
	GetMessage() []*chan *MessageResponse
	GetMetadata() []*chan *MetadataResponse
	GetSpeechStarted() []*chan *SpeechStartedResponse
	GetUtteranceEnd() []*chan *UtteranceEndResponse
	GetClose() []*chan *CloseResponse
	GetError() []*chan *ErrorResponse
	GetUnhandled() []*chan *[]byte
}

/*
Callback Interfaces
*/
// LiveMessageCallback is a callback used to receive notifcations for platforms messages
type LiveMessageCallback interface {
	Open(or *OpenResponse) error
	Message(mr *MessageResponse) error
	Metadata(md *MetadataResponse) error
	SpeechStarted(ssr *SpeechStartedResponse) error
	UtteranceEnd(ur *UtteranceEndResponse) error
	Close(cr *CloseResponse) error
	Error(er *ErrorResponse) error
	UnhandledEvent(byData []byte) error
}
