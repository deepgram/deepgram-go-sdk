// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

// DefaultCallbackHandler is a default callback handler for live transcription
// Simply prints the transcript to stdout
type DefaultCallbackHandler struct{}

// NewDefaultCallbackHandler creates a new DefaultCallbackHandler
func NewDefaultCallbackHandler() DefaultCallbackHandler {
	return DefaultCallbackHandler{}
}

// Message is the callback for a message
func (dch DefaultCallbackHandler) Message(mr *interfaces.MessageResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(mr)
		if err != nil {
			klog.V(1).Infof("Message json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nMessage Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		klog.V(7).Infof("DEEPGRAM - no transcript")
		return nil
	}
	fmt.Printf("\n%s\n", sentence)

	return nil
}

// Metadata is the callback for a metadata
func (dch DefaultCallbackHandler) Metadata(md *interfaces.MetadataResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(md)
		if err != nil {
			klog.V(1).Infof("Metadata json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nMetadata Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	fmt.Printf("\nMetadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	fmt.Printf("Metadata.Channels: %d\n", md.Channels)
	fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))

	return nil
}

func (dch DefaultCallbackHandler) UtteranceEnd() error {
	fmt.Printf("\nUtteranceEnd \n")
	return nil
}

func (dch DefaultCallbackHandler) Error(er *interfaces.ErrorResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(er)
		if err != nil {
			klog.V(1).Infof("Error json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nError Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	fmt.Printf("\nError.Type: %s\n", er.Type)
	fmt.Printf("Error.Message: %s\n", er.Message)
	fmt.Printf("Error.Description: %s\n\n", er.Description)

	return nil
}
