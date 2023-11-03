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

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

type DefaultCallbackHandler struct {
	sb strings.Builder
}

func (dch DefaultCallbackHandler) Message(mr *interfaces.MessageResponse) error {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}

	if strings.Compare(strings.ToLower(debugStr), "true") == 0 {
		data, err := json.Marshal(mr)
		if err != nil {
			klog.V(1).Infof("RecognitionResult json.Marshal failed. Err: %v\n", err)
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
	fmt.Printf("%s\n", sentence)

	return nil
}

func NewDefaultCallbackHandler() DefaultCallbackHandler {
	return DefaultCallbackHandler{}
}
