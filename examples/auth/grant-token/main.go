// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"fmt"
	"os"

	authAPI "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1"
	authInterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/auth/v1/interfaces"
	authClient "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/auth"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v1"
)

func main() {
	// Initialize library
	authClient.Init(authClient.InitLib{
		LogLevel: authClient.LogLevelTrace,
	})

	ctx := context.Background()

	fmt.Printf("=== Deepgram Dual Authentication Demo ===\n\n")

	// ====================================
	// Phase 1: API Key Authentication
	// ====================================
	fmt.Printf("1. Using API Key Authentication (Token)\n")

	// Create auth client with API key (will use "Authorization: token <API_KEY>")
	dgAuth := authClient.NewWithDefaults()
	authAPIClient := authAPI.New(dgAuth)

	// Grant a temporary access token using API key
	fmt.Printf("   Requesting access token using API key...\n")
	respToken, err := authAPIClient.GrantToken(ctx, nil)
	if err != nil {
		fmt.Printf("   ❌ GrantToken failed. Err: %v\n", err)
		os.Exit(1)
	}

	accessToken := respToken.AccessToken
	ttl := respToken.ExpiresIn
	fmt.Printf("   ✅ Access token granted successfully!\n")
	fmt.Printf("   Token: %s\n", accessToken[:16]+"...")
	fmt.Printf("   TTL: %.0f seconds\n\n", ttl)

	// ====================================
	// Phase 1.5: API Key Authentication with custom TTL
	// ====================================
	fmt.Printf("1.5. Using API Key Authentication with custom TTL (60 seconds)\n")

	// Grant a temporary access token with custom TTL
	fmt.Printf("   Requesting access token with 60-second TTL...\n")
	customTTL := 60
	respTokenCustom, err := authAPIClient.GrantToken(ctx, &authInterfaces.GrantTokenRequest{
		TTLSeconds: &customTTL,
	})
	if err != nil {
		fmt.Printf("   ❌ GrantToken with custom TTL failed. Err: %v\n", err)
		os.Exit(1)
	}

	accessTokenCustom := respTokenCustom.AccessToken
	ttlCustom := respTokenCustom.ExpiresIn
	fmt.Printf("   ✅ Access token with custom TTL granted successfully!\n")
	fmt.Printf("   Token: %s\n", accessTokenCustom[:16]+"...")
	fmt.Printf("   TTL: %.0f seconds\n\n", ttlCustom)

	// ====================================
	// Phase 2: Dual Authentication Features
	// ====================================
	fmt.Printf("2. Demonstrating Dual Authentication Features\n")

	// Create client options with access token
	bearerOptions := &interfaces.ClientOptions{
		AccessToken: accessToken, // Priority 1: Explicit AccessToken
	}

	fmt.Printf("   Created ClientOptions with AccessToken\n")
	token, isBearer := bearerOptions.GetAuthToken()
	if isBearer {
		fmt.Printf("   ✅ Priority resolution: Using Bearer token (%s...)\n", token[:16])
	} else {
		fmt.Printf("   ❌ Expected Bearer token but got API key\n")
	}

	// Test priority system with both credentials
	dualOptions := &interfaces.ClientOptions{
		APIKey:      os.Getenv("DEEPGRAM_API_KEY"), // Lower priority
		AccessToken: accessToken,                   // Higher priority
	}

	fmt.Printf("\n   Testing priority with both credentials set:\n")
	token, isBearer = dualOptions.GetAuthToken()
	if isBearer {
		fmt.Printf("   ✅ AccessToken takes priority (Bearer: %s...)\n", token[:16])
	} else {
		fmt.Printf("   ❌ Expected AccessToken priority\n")
	}

	// ====================================
	// Phase 3: Dynamic Authentication Switching
	// ====================================
	fmt.Printf("\n3. Dynamic Authentication Switching\n")

	// Switch to API key authentication dynamically
	fmt.Printf("   Clearing AccessToken to switch to APIKey...\n")
	dualOptions.SetAccessToken("")
	token, isBearer = dualOptions.GetAuthToken()
	if !isBearer && token != "" {
		fmt.Printf("   ✅ Switched to Token auth (APIKey: %s...)\n", token[:16])
	} else {
		fmt.Printf("   ❌ Failed to switch to API key\n")
	}

	// Switch back to Bearer token
	fmt.Printf("   Switching back to Bearer token...\n")
	dualOptions.SetAccessToken(accessToken)
	token, isBearer = dualOptions.GetAuthToken()
	if isBearer {
		fmt.Printf("   ✅ Switched back to Bearer auth (%s...)\n", token[:16])
	} else {
		fmt.Printf("   ❌ Failed to switch back to Bearer\n")
	}

	// Test APIKey update
	fmt.Printf("   Testing API key update...\n")
	dualOptions.SetAPIKey("test_new_api_key")
	dualOptions.SetAccessToken("") // Clear access token to use API key
	token, isBearer = dualOptions.GetAuthToken()
	if !isBearer && token == "test_new_api_key" {
		fmt.Printf("   ✅ API key updated successfully\n")
	} else {
		fmt.Printf("   ❌ API key update failed\n")
	}

	// ====================================
	// Phase 4: Environment Variable Priority Demo
	// ====================================
	fmt.Printf("\n4. Environment Variable Priority\n")
	fmt.Printf("   Priority order:\n")
	fmt.Printf("   1. Explicit AccessToken parameter (highest)\n")
	fmt.Printf("   2. Explicit APIKey parameter\n")
	fmt.Printf("   3. DEEPGRAM_ACCESS_TOKEN environment variable\n")
	fmt.Printf("   4. DEEPGRAM_API_KEY environment variable (lowest)\n\n")

	// Test environment variable resolution
	envOptions := &interfaces.ClientOptions{}
	err = envOptions.Parse()
	if err != nil {
		fmt.Printf("   ❌ Environment parsing failed: %v\n", err)
	} else {
		token, isBearer := envOptions.GetAuthToken()
		if token != "" {
			if isBearer {
				fmt.Printf("   Environment resolved: Bearer token (%s...)\n", token[:16])
			} else {
				fmt.Printf("   Environment resolved: API key (%s...)\n", token[:16])
			}
		} else {
			fmt.Printf("   No credentials found in environment\n")
		}
	}

	// Show current environment variable status
	if os.Getenv("DEEPGRAM_ACCESS_TOKEN") != "" {
		fmt.Printf("   DEEPGRAM_ACCESS_TOKEN: Set (%s...)\n", os.Getenv("DEEPGRAM_ACCESS_TOKEN")[:16])
	} else {
		fmt.Printf("   DEEPGRAM_ACCESS_TOKEN: Not set\n")
	}

	if os.Getenv("DEEPGRAM_API_KEY") != "" {
		fmt.Printf("   DEEPGRAM_API_KEY: Set (%s...)\n", os.Getenv("DEEPGRAM_API_KEY")[:16])
	} else {
		fmt.Printf("   DEEPGRAM_API_KEY: Not set\n")
	}

	// ====================================
	// Phase 5: Backward Compatibility Verification
	// ====================================
	fmt.Printf("\n5. Backward Compatibility Verification\n")

	// Test that existing API key only usage still works
	legacyOptions := &interfaces.ClientOptions{
		APIKey: os.Getenv("DEEPGRAM_API_KEY"),
	}
	token, isBearer = legacyOptions.GetAuthToken()
	if !isBearer && token != "" {
		fmt.Printf("   ✅ Legacy API key usage: Compatible\n")
	} else {
		fmt.Printf("   ❌ Legacy API key usage: Broken\n")
	}

	// Test environment variable fallback (existing behavior)
	emptyOptions := &interfaces.ClientOptions{}
	err = emptyOptions.Parse()
	if err == nil {
		fmt.Printf("   ✅ Environment fallback: Compatible\n")
	} else {
		fmt.Printf("   ❌ Environment fallback: Broken - %v\n", err)
	}

	fmt.Printf("\n✅ Dual authentication demo completed successfully!\n")
	fmt.Printf("\nKey Features Demonstrated:\n")
	fmt.Printf("• API Key → Bearer Token workflow\n")
	fmt.Printf("• Priority-based credential resolution\n")
	fmt.Printf("• Dynamic authentication switching\n")
	fmt.Printf("• Environment variable support\n")
	fmt.Printf("• Backward compatibility\n")
}
