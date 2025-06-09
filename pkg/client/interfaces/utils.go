// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package interfaces

import (
	"context"
	"net/http"

	interfacesv1 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

// DgAgent is the agent version
var DgAgent = interfacesv1.DgAgent

// signer
type Signer = interfacesv1.Signer
type SignerContext = interfacesv1.SignerContext

func WithSigner(ctx context.Context, s Signer) context.Context {
	return interfacesv1.WithSigner(ctx, s)
}

// headers
type HeadersContext = interfacesv1.HeadersContext

func WithCustomHeaders(ctx context.Context, headers http.Header) context.Context {
	return interfacesv1.WithCustomHeaders(ctx, headers)
}

// parameters
type ParametersContext = interfacesv1.ParametersContext

func WithCustomParameters(ctx context.Context, params map[string][]string) context.Context {
	return interfacesv1.WithCustomParameters(ctx, params)
}

// common structs found throughout the SDK
type RawResponse = interfacesv1.RawResponse
type DeepgramWarning = interfacesv1.DeepgramWarning
type DeepgramError = interfacesv1.DeepgramError
type StatusError = interfacesv1.StatusError
