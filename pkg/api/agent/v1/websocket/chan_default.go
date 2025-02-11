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

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/agent/v1/websocket/interfaces"
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
		debugWebsocket:               strings.EqualFold(debugStr, "true"),
		debugWebsocketVerbose:        strings.EqualFold(debugExtStr, "true"),
		binaryChan:                   make(chan *[]byte),
		openChan:                     make(chan *interfaces.OpenResponse),
		welcomeResponse:              make(chan *interfaces.WelcomeResponse),
		conversationTextResponse:     make(chan *interfaces.ConversationTextResponse),
		userStartedSpeakingResponse:  make(chan *interfaces.UserStartedSpeakingResponse),
		agentThinkingResponse:        make(chan *interfaces.AgentThinkingResponse),
		functionCallRequestResponse:  make(chan *interfaces.FunctionCallRequestResponse),
		functionCallingResponse:      make(chan *interfaces.FunctionCallingResponse),
		agentStartedSpeakingResponse: make(chan *interfaces.AgentStartedSpeakingResponse),
		agentAudioDoneResponse:       make(chan *interfaces.AgentAudioDoneResponse),
		injectionRefusedResponse:     make(chan *interfaces.InjectionRefusedResponse),
		keepAliveResponse:            make(chan *interfaces.KeepAlive),
		closeChan:                    make(chan *interfaces.CloseResponse),
		errorChan:                    make(chan *interfaces.ErrorResponse),
		unhandledChan:                make(chan *[]byte),
	}

	go func() {
		err := handler.Run()
		if err != nil {
			klog.V(1).Infof("handler.Run failed. Err: %v\n", err)
		}
	}()

	return handler
}

// GetBinary returns the binary channels
func (dch DefaultChanHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&dch.binaryChan}
}

// GetOpen returns the open channels
func (dch DefaultChanHandler) GetOpen() []*chan *interfaces.OpenResponse {
	return []*chan *interfaces.OpenResponse{&dch.openChan}
}

// GetWelcomeResponse returns the welcome response channels
func (dch DefaultChanHandler) GetWelcome() []*chan *interfaces.WelcomeResponse {
	return []*chan *interfaces.WelcomeResponse{&dch.welcomeResponse}
}

// GetConversationTextResponse returns the conversation text response channels
func (dch DefaultChanHandler) GetConversationText() []*chan *interfaces.ConversationTextResponse {
	return []*chan *interfaces.ConversationTextResponse{&dch.conversationTextResponse}
}

// GetUserStartedSpeakingResponse returns the user started speaking response channels
func (dch DefaultChanHandler) GetUserStartedSpeaking() []*chan *interfaces.UserStartedSpeakingResponse {
	return []*chan *interfaces.UserStartedSpeakingResponse{&dch.userStartedSpeakingResponse}
}

// GetAgentThinkingResponse returns the agent thinking response channels
func (dch DefaultChanHandler) GetAgentThinking() []*chan *interfaces.AgentThinkingResponse {
	return []*chan *interfaces.AgentThinkingResponse{&dch.agentThinkingResponse}
}

// GetFunctionCallRequestResponse returns the function call request response channels
func (dch DefaultChanHandler) GetFunctionCallRequest() []*chan *interfaces.FunctionCallRequestResponse {
	return []*chan *interfaces.FunctionCallRequestResponse{&dch.functionCallRequestResponse}
}

// GetFunctionCallingResponse returns the function calling response channels
func (dch DefaultChanHandler) GetFunctionCalling() []*chan *interfaces.FunctionCallingResponse {
	return []*chan *interfaces.FunctionCallingResponse{&dch.functionCallingResponse}
}

// GetAgentStartedSpeakingResponse returns the agent started speaking response channels
func (dch DefaultChanHandler) GetAgentStartedSpeaking() []*chan *interfaces.AgentStartedSpeakingResponse {
	return []*chan *interfaces.AgentStartedSpeakingResponse{&dch.agentStartedSpeakingResponse}
}

// GetAgentAudioDoneResponse returns the agent audio done response channels
func (dch DefaultChanHandler) GetAgentAudioDone() []*chan *interfaces.AgentAudioDoneResponse {
	return []*chan *interfaces.AgentAudioDoneResponse{&dch.agentAudioDoneResponse}
}

// GetClose returns the close channels
func (dch DefaultChanHandler) GetClose() []*chan *interfaces.CloseResponse {
	return []*chan *interfaces.CloseResponse{&dch.closeChan}
}

// GetError returns the error channels
func (dch DefaultChanHandler) GetError() []*chan *interfaces.ErrorResponse {
	return []*chan *interfaces.ErrorResponse{&dch.errorChan}
}

// GetInjectionRefused returns the injection refused channels
func (dch DefaultChanHandler) GetInjectionRefused() []*chan *interfaces.InjectionRefusedResponse {
	return []*chan *interfaces.InjectionRefusedResponse{&dch.injectionRefusedResponse}
}

// GetKeepAlive returns the keep alive channels
func (dch DefaultChanHandler) GetKeepAlive() []*chan *interfaces.KeepAlive {
	return []*chan *interfaces.KeepAlive{&dch.keepAliveResponse}
}

// GetUnhandled returns the unhandled event channels
func (dch DefaultChanHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

// GetSettingsApplied returns the settings applied response channels
func (dch DefaultChanHandler) GetSettingsApplied() []*chan *interfaces.SettingsAppliedResponse {
	return []*chan *interfaces.SettingsAppliedResponse{&dch.settingsAppliedResponse}
}

// Open is the callback for when the connection opens
//
//nolint:funlen,gocyclo // this is a complex function. keep as is
func (dch DefaultChanHandler) Run() error {
	wgReceivers := sync.WaitGroup{}

	// binary channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for br := range dch.binaryChan {
			fmt.Printf("\n\n[Binary Data]\n\n")
			fmt.Printf("Size: %d\n\n", len(*br))

			if dch.debugWebsocket {
				fmt.Printf("Hex Dump: %x...\n\n", (*br)[:20])
			}
			if dch.debugWebsocketVerbose {
				fmt.Printf("Dumping to verbose.wav\n")
				file, err := os.OpenFile("verbose.wav", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
				if err != nil {
					fmt.Printf("Failed to open file. Err: %v\n", err)
					continue
				}

				_, err = file.Write(*br)
				file.Close()

				if err != nil {
					fmt.Printf("Failed to write to file. Err: %v\n", err)
					continue
				}
			}
		}
	}()

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

	// welcome response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for wr := range dch.welcomeResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(wr)
				if err != nil {
					klog.V(1).Infof("Welcome json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nWelcome Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[WelcomeResponse]\n\n")
		}
	}()

	// conversation text response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ctr := range dch.conversationTextResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(ctr)
				if err != nil {
					klog.V(1).Infof("ConversationText json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nConversationText Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[ConversationTextResponse]\n\n")
		}
	}()

	// user started speaking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ussr := range dch.userStartedSpeakingResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(ussr)
				if err != nil {
					klog.V(1).Infof("UserStartedSpeaking json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nUserStartedSpeaking Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[UserStartedSpeakingResponse]\n\n")
		}
	}()

	// agent thinking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for atr := range dch.agentThinkingResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(atr)
				if err != nil {
					klog.V(1).Infof("AgentThinking json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nAgentThinking Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[AgentThinkingResponse]\n\n")
		}
	}()

	// function call request response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for fcrr := range dch.functionCallRequestResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(fcrr)
				if err != nil {
					klog.V(1).Infof("FunctionCallRequest json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nFunctionCallRequest Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[FunctionCallRequestResponse]\n\n")
		}
	}()

	// function calling response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for fcr := range dch.functionCallingResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(fcr)
				if err != nil {
					klog.V(1).Infof("FunctionCalling json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nFunctionCalling Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[FunctionCallingResponse]\n\n")
		}
	}()

	// agent started speaking response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for assr := range dch.agentStartedSpeakingResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(assr)
				if err != nil {
					klog.V(1).Infof("AgentStartedSpeaking json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nAgentStartedSpeaking Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[AgentStartedSpeakingResponse]\n\n")
		}
	}()

	// agent audio done response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for aadr := range dch.agentAudioDoneResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(aadr)
				if err != nil {
					klog.V(1).Infof("AgentAudioDone json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nAgentAudioDone Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[AgentAudioDoneResponse]\n\n")
		}
	}()

	// keep alive response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ka := range dch.keepAliveResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(ka)
				if err != nil {
					klog.V(1).Infof("KeepAlive json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nKeepAlive Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[KeepAliveResponse]\n\n")
		}
	}()

	// settings applied response channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for sa := range dch.settingsAppliedResponse {
			if dch.debugWebsocket {
				data, err := json.Marshal(sa)
				if err != nil {
					klog.V(1).Infof("SettingsApplied json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nSettingsApplied Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n\n[SettingsAppliedResponse]\n\n")
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
