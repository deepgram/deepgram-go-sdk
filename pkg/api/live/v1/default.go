// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package live

import (

	// prettyjson "github.com/hokaccha/go-prettyjson"

	"log"
	"strings"

	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
)

type DefaultCallbackHandler struct {
	sb strings.Builder
}

func (dch DefaultCallbackHandler) Message(mr *interfaces.MessageResponse) error {
	// DEBUG: just print the transcription
	// data, err := json.Marshal(mr)
	// if err != nil {
	// 	log.Printf("RecognitionResult json.Marshal failed. Err: %v\n", err)
	// 	return err
	// }

	// prettyJson, err := prettyjson.Format(data)
	// if err != nil {
	// 	log.Printf("prettyjson.Marshal failed. Err: %v\n", err)
	// 	return err
	// }
	// log.Printf("\n\nMessage Object:\n%s\n\n", prettyJson)

	// Only print the final transcript
	// // handle the message
	// sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	// if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
	// 	// klog.V(7).Infof("DEEPGRAM - no transcript")
	// 	return nil
	// }

	// isFinal := mr.SpeechFinal
	// sentence = strings.ToLower(sentence)
	// dch.sb.WriteString(sentence)

	// // // debug
	// // klog.V(4).Infof("transcription result: text = %s, final = %t", i.sb.String(), isFinal)

	// if !isFinal {
	// 	// klog.V(7).Infof("DEEPGRAM - not final")
	// 	return nil
	// }

	// // debug
	// log.Printf("Deepgram transcription: text = %s, final = %t", dch.sb.String(), isFinal)
	// dch.sb.Reset()

	// handle the message
	sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

	if len(mr.Channel.Alternatives) == 0 || len(sentence) == 0 {
		// klog.V(7).Infof("DEEPGRAM - no transcript")
		return nil
	}
	log.Printf("%s\n", sentence)

	return nil
}

func NewDefaultCallbackHandler() DefaultCallbackHandler {
	return DefaultCallbackHandler{}
}
