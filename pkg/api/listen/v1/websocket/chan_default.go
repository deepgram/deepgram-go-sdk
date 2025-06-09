// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package websocketv1

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	prettyjson "github.com/hokaccha/go-prettyjson"
	klog "k8s.io/klog/v2"

	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/websocket/interfaces"
)

// NewDefaultChanHandler creates a new DefaultChanHandler
func NewDefaultChanHandler() *DefaultChanHandler {
	var debugStr string
	if v := os.Getenv("DEEPGRAM_DEBUG"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG found")
		debugStr = v
	}
	var debugExtStr string
	if v := os.Getenv("DEEPGRAM_DEBUG_VERBOSE"); v != "" {
		klog.V(4).Infof("DEEPGRAM_DEBUG_VERBOSE found")
		debugExtStr = v
	}
	handler := &DefaultChanHandler{
		debugWebsocket:        strings.EqualFold(debugStr, "true"),
		debugWebsocketVerbose: strings.EqualFold(debugExtStr, "true"),
		openChan:              make(chan *interfaces.OpenResponse),
		messageChan:           make(chan *interfaces.MessageResponse),
		metadataChan:          make(chan *interfaces.MetadataResponse),
		speechStartedChan:     make(chan *interfaces.SpeechStartedResponse),
		utteranceEndChan:      make(chan *interfaces.UtteranceEndResponse),
		closeChan:             make(chan *interfaces.CloseResponse),
		errorChan:             make(chan *interfaces.ErrorResponse),
		unhandledChan:         make(chan *[]byte),
	}

	go func() {
		err := handler.Run()
		if err != nil {
			klog.V(1).Infof("handler.Run failed. Err: %v\n", err)
		}
	}()

	return handler
}

// GetOpen returns the open channels
func (dch DefaultChanHandler) GetOpen() []*chan *interfaces.OpenResponse {
	return []*chan *interfaces.OpenResponse{&dch.openChan}
}

// GetMessage returns the message channels
func (dch DefaultChanHandler) GetMessage() []*chan *interfaces.MessageResponse {
	return []*chan *interfaces.MessageResponse{&dch.messageChan}
}

// GetMetadata returns the metadata channels
func (dch DefaultChanHandler) GetMetadata() []*chan *interfaces.MetadataResponse {
	return []*chan *interfaces.MetadataResponse{&dch.metadataChan}
}

// GetSpeechStarted returns the speech started channels
func (dch DefaultChanHandler) GetSpeechStarted() []*chan *interfaces.SpeechStartedResponse {
	return []*chan *interfaces.SpeechStartedResponse{&dch.speechStartedChan}
}

// GetUtteranceEnd returns the utterance end channels
func (dch DefaultChanHandler) GetUtteranceEnd() []*chan *interfaces.UtteranceEndResponse {
	return []*chan *interfaces.UtteranceEndResponse{&dch.utteranceEndChan}
}

// GetClose returns the close channels
func (dch DefaultChanHandler) GetClose() []*chan *interfaces.CloseResponse {
	return []*chan *interfaces.CloseResponse{&dch.closeChan}
}

// GetError returns the error channels
func (dch DefaultChanHandler) GetError() []*chan *interfaces.ErrorResponse {
	return []*chan *interfaces.ErrorResponse{&dch.errorChan}
}

// GetUnhandled returns the unhandled event channels
func (dch DefaultChanHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

// Open is the callback for when the connection opens
//
//nolint:funlen,gocyclo // this is a complex function. keep as is
func (dch DefaultChanHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// open channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for or := range dch.openChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(or)
				if err != nil {
					klog.V(1).Infof("Open json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nOpen Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[OpenResponse]\n\n")
		}
	}()

	// message channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.messageChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(mr)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nMessage Object:\n%s\n\n", prettyJSON)
			}

			sentence := strings.TrimSpace(mr.Channel.Alternatives[0].Transcript)

			if len(mr.Channel.Alternatives) == 0 || sentence == "" {
				klog.V(7).Infof("DEEPGRAM - no transcript")
				continue
			}

			if mr.IsFinal {
				fmt.Printf("\n[MessageResponse] (Final) %s\n", sentence)
			} else {
				fmt.Printf("\n[MessageResponse] (Interim) %s\n", sentence)
			}
		}
	}()

	// metadata channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.metadataChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(mr)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nMetadata Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\nMetadata.RequestID: %s\n", strings.TrimSpace(mr.RequestID))
			fmt.Printf("Metadata.Channels: %d\n", mr.Channels)
			fmt.Printf("Metadata.Created: %s\n\n", strings.TrimSpace(mr.Created))
		}
	}()

	// speech started channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ssr := range dch.speechStartedChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(ssr)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nSpeechStarted Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[SpeechStarted]\n")
			if dch.debugWebsocketVerbose {
				fmt.Printf("\n\nSpeechStarted.Timestamp: %f\n", ssr.Timestamp)
				fmt.Printf("SpeechStarted.Channels:\n")
				for _, val := range ssr.Channel {
					fmt.Printf("\tChannel: %d\n", val)
				}
				fmt.Printf("\n")
			}
		}
	}()

	// utterance end channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for uer := range dch.utteranceEndChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(uer)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nUtteranceEnd Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[UtteranceEnd]\n")
			if dch.debugWebsocketVerbose {
				fmt.Printf("\nUtteranceEnd.Timestamp: %f\n", uer.LastWordEnd)
				fmt.Printf("UtteranceEnd.Channel: %d\n\n", uer.Channel)
			}
		}
	}()

	// close channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for cr := range dch.closeChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(cr)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nClose Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[CloseResponse]\n\n")
		}
	}()

	// error channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for er := range dch.errorChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(er)
				if err != nil {
					klog.V(1).Infof("Close json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nError Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[ErrorResponse]\n")
			fmt.Printf("\nError.Type: %s\n", er.ErrCode)
			fmt.Printf("Error.Message: %s\n", er.ErrMsg)
			fmt.Printf("Error.Description: %s\n\n", er.Description)
			fmt.Printf("Error.Variant: %s\n\n", er.Variant)
		}
	}()

	// unhandled event channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for byData := range dch.unhandledChan {
			if dch.debugWebsocket {
				prettyJSON, err := prettyjson.Format(*byData)
				if err != nil {
					klog.V(2).Infof("\n\nRaw Data:\n%s\n\n", string(*byData))
				} else {
					klog.V(2).Infof("\n\nError Object:\n%s\n\n", prettyJSON)
				}
			}

			fmt.Printf("\n[UnhandledEvent]")
			fmt.Printf("Dump:\n%s\n\n", string(*byData))
		}
	}()

	// wait for all receivers to finish
	wgReceivers.Wait()

	return nil
}
