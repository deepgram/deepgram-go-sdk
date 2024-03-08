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

// Open is the callback for when the connection opens
func (dch DefaultCallbackHandler) Open(or *interfaces.OpenResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(or)
		if err != nil {
			klog.V(1).Infof("Open json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nOpen Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	fmt.Printf("\n\nOpenResponse\n\n")

	return nil
}

// Message is the callback for a transcription message
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

// Metadata is the callback for information about the connection
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
	fmt.Printf("\n\nMetadata.RequestID: %s\n", strings.TrimSpace(md.RequestID))
	fmt.Printf("Metadata.Channels: %d\n", md.Channels)
	fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(md.Created))

	return nil
}

// SpeechStarted is when VAD detects noise
func (dch DefaultCallbackHandler) SpeechStarted(ssr *interfaces.SpeechStartedResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(ssr)
		if err != nil {
			klog.V(1).Infof("SpeechStarted json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nSpeechStarted Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	fmt.Printf("\n\nSpeechStarted.Timestamp: %f\n", ssr.Timestamp)
	fmt.Printf("SpeechStarted.Channels:\n")
	for _, val := range ssr.Channel {
		fmt.Printf("\tChannel: %d\n", val)
	}
	fmt.Printf("\n")

	return nil
}

// UtteranceEnd is the callback for when a channel goes silent
func (dch DefaultCallbackHandler) UtteranceEnd(ur *interfaces.UtteranceEndResponse) error {
	fmt.Printf("\nUtteranceEnd.Timestamp: %f\n", ur.LastWordEnd)
	fmt.Printf("UtteranceEnd.Channel: %d\n\n", ur.Channel)
	return nil
}

// Close is the callback for when the connection closes
func (dch DefaultCallbackHandler) Close(or *interfaces.CloseResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(or)
		if err != nil {
			klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
			return err
		}

		prettyJson, err := prettyjson.Format(data)
		if err != nil {
			klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
			return err
		}
		klog.V(2).Infof("\n\nClose Object:\n%s\n\n", prettyJson)

		return nil
	}

	// handle the message
	fmt.Printf("\n\nCloseResponse\n\n")

	return nil
}

// Error is the callback for a error messages
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

// UnhandledEvent is the callback for unknown messages
func (dch DefaultCallbackHandler) UnhandledEvent(byData []byte) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		prettyJson, err := prettyjson.Format(byData)
		if err != nil {
			klog.V(2).Infof("\n\nRaw Data:\n%s\n\n", string(byData))
		} else {
			klog.V(2).Infof("\n\nError Object:\n%s\n\n", prettyJson)
		}

		return nil
	}

	// handle the message
	fmt.Printf("\n\nUnhandledEvent\n%s\n\n", string(byData))

	return nil
}
