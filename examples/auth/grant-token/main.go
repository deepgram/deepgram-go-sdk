// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth"
)

func main() {
	// init library
	client.Init(client.InitLib{
		LogLevel: client.LogLevelTrace, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	// context
	ctx := context.Background()

	//client
	dg := client.NewWithDefaults()
	authClient := api.New(dg)

	// list projects
	respToken, err := authClient.GrantToken(ctx)
	if err != nil {
		fmt.Printf("GrantToken failed. Err: %v\n", err)
		os.Exit(1)
	}

	var token = respToken.AccessToken
	var ttl = respToken.ExpiresIn

	fmt.Printf("GrantToken() - Token: %s, TTL: %f\n", token, ttl)
}
