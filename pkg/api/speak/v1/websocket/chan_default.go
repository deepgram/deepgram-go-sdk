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

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/websocket/interfaces"
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

	handler := DefaultChanHandler{
		debugWebsocket:        strings.EqualFold(debugStr, "true"),
		debugWebsocketVerbose: strings.EqualFold(debugExtStr, "true"),
		binaryChan:            make(chan *[]byte),
		openChan:              make(chan *interfaces.OpenResponse),
		metadataChan:          make(chan *interfaces.MetadataResponse),
		flushedChan:           make(chan *interfaces.FlushedResponse),
		clearedChan:           make(chan *interfaces.ClearedResponse),
		closeChan:             make(chan *interfaces.CloseResponse),
		warningChan:           make(chan *interfaces.WarningResponse),
		errorChan:             make(chan *interfaces.ErrorResponse),
		unhandledChan:         make(chan *[]byte),
	}

	go func() {
		err := handler.Run()
		if err != nil {
			klog.V(1).Infof("handler.Run failed. Err: %v\n", err)
		}
	}()

	return &handler
}

// GetBinary returns the binary event channels
func (dch *DefaultChanHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&dch.binaryChan}
}

// GetOpen returns the open channels
func (dch *DefaultChanHandler) GetOpen() []*chan *interfaces.OpenResponse {
	return []*chan *interfaces.OpenResponse{&dch.openChan}
}

// GetMetadata returns the metadata channels
func (dch *DefaultChanHandler) GetMetadata() []*chan *interfaces.MetadataResponse {
	return []*chan *interfaces.MetadataResponse{&dch.metadataChan}
}

// GetfFlush returns the flush channels
func (dch *DefaultChanHandler) GetFlush() []*chan *interfaces.FlushedResponse {
	return []*chan *interfaces.FlushedResponse{&dch.flushedChan}
}

// GetCleared returns the cleared channels
func (dch *DefaultChanHandler) GetClear() []*chan *interfaces.ClearedResponse {
	return []*chan *interfaces.ClearedResponse{&dch.clearedChan}
}

// GetClose returns the close channels
func (dch *DefaultChanHandler) GetClose() []*chan *interfaces.CloseResponse {
	return []*chan *interfaces.CloseResponse{&dch.closeChan}
}

// GetWarning returns the warning channels
func (dch *DefaultChanHandler) GetWarning() []*chan *interfaces.WarningResponse {
	return []*chan *interfaces.WarningResponse{&dch.warningChan}
}

// GetError returns the error channels
func (dch *DefaultChanHandler) GetError() []*chan *interfaces.ErrorResponse {
	return []*chan *interfaces.ErrorResponse{&dch.errorChan}
}

// GetUnhandled returns the unhandled event channels
func (dch *DefaultChanHandler) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&dch.unhandledChan}
}

// Open is the callback for when the connection opens
//
//nolint:funlen,gocyclo // this is a complex function. keep as is
func (dch *DefaultChanHandler) Run() error {
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

	// metadata channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for mr := range dch.metadataChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(mr)
				if err != nil {
					klog.V(1).Infof("Metadata json.Marshal failed. Err: %v\n", err)
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
		}
	}()

	// flushed channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ssr := range dch.flushedChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(ssr)
				if err != nil {
					klog.V(1).Infof("Flush json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nFlushed Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[Flushed]\n")
		}
	}()

	// cleared channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for ssr := range dch.clearedChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(ssr)
				if err != nil {
					klog.V(1).Infof("Clear json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nCleared Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[Cleared]\n")
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

	// warning channel
	wgReceivers.Add(1)
	go func() {
		defer wgReceivers.Done()

		for wr := range dch.warningChan {
			if dch.debugWebsocket {
				data, err := json.Marshal(wr)
				if err != nil {
					klog.V(1).Infof("Warning json.Marshal failed. Err: %v\n", err)
					continue
				}

				prettyJSON, err := prettyjson.Format(data)
				if err != nil {
					klog.V(1).Infof("prettyjson.Marshal failed. Err: %v\n", err)
					continue
				}
				klog.V(2).Infof("\n\nWarning Object:\n%s\n\n", prettyJSON)
			}

			fmt.Printf("\n[Warning]\n")
			fmt.Printf("\nWarning.Type: %s\n", wr.WarnCode)
			fmt.Printf("Warning.Message: %s\n", wr.WarnMsg)
			fmt.Printf("Warning.Description: %s\n\n", wr.Description)
			fmt.Printf("Warning.Variant: %s\n\n", wr.Variant)
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
					klog.V(1).Infof("Error json.Marshal failed. Err: %v\n", err)
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
