// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package commonv2

import (
	commonv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/common/v1"
)

// external constants
const (
	DefaultConnectRetry = commonv1.DefaultConnectRetry

	ChunkSize        = commonv1.ChunkSize
	TerminationSleep = commonv1.TerminationSleep

	// socket errors
	FatalReadSocketErr  = commonv1.FatalReadSocketErr
	FatalWriteSocketErr = commonv1.FatalWriteSocketErr
	UseOfClosedSocket   = commonv1.UseOfClosedSocket
	UnknownDeepgramErr  = commonv1.UnknownDeepgramErr

	// socket successful close error
	SuccessfulSocketErr = commonv1.SuccessfulSocketErr
)

// errors
var (
	ErrInvalidInput        = commonv1.ErrInvalidInput
	ErrInvalidConnection   = commonv1.ErrInvalidConnection
	ErrFatalPanicRecovered = commonv1.ErrFatalPanicRecovered
)
