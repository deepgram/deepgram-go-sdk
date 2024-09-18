// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfacesv1

// These are the message types that can be received from the live API
type TypeResponse string

const (
	// message types
	TypeOpenResponse  TypeResponse = "Open"
	TypeCloseResponse TypeResponse = "Close"
	TypeErrorResponse TypeResponse = "Error"
)
