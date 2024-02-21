// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

// This package defines the Pre-recorded API for Deepgram
package speak

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"strconv"

	klog "k8s.io/klog/v2"

	api "github.com/deepgram/deepgram-go-sdk/pkg/api/speak/v1/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/speak"
)

type SpeakClient struct {
	*client.Client
}

func New(client *client.Client) *SpeakClient {
	return &SpeakClient{client}
}

func (c *SpeakClient) ToStream(ctx context.Context, text string, options interfaces.SpeakOptions) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speak.ToStream ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}
	err := options.Check()
	if err != nil {
		klog.V(1).Infof("SpeakOptions.Check() failed. Err: %v\n", err)
		return nil, err
	}

	// send the file!
	retVal := make(map[string]string)
	retVal["content-type"] = ""
	retVal["request-id"] = ""
	retVal["model-uuid"] = ""
	retVal["model-name"] = ""
	retVal["char-count"] = ""
	retVal["transfer-encoding"] = ""
	retVal["date"] = ""
	var resp bytes.Buffer

	err = c.Client.DoText(ctx, text, options, &retVal, &resp)
	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				klog.V(6).Infof("speak.ToStream LEAVE\n")
				return nil, err
			}
		}

		klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		klog.V(6).Infof("speak.ToStream LEAVE\n")
		return nil, err
	}

	charCnt, err := strconv.Atoi(retVal["char-count"])
	if err != nil {
		klog.V(1).Infof("strconv.Atoi failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToStream LEAVE\n")
		return nil, err
	}

	var result api.SpeakResponse
	result.ContextType = retVal["content-type"]
	result.RequestId = retVal["request-id"]
	result.ModelUuid = retVal["model-uuid"]
	result.ModelName = retVal["model-name"]
	result.Characters = charCnt
	result.TransferEncoding = retVal["transfer-encoding"]
	result.Date = retVal["date"]
	result.Stream = &resp

	klog.V(3).Infof("Speak Succeeded\n")
	klog.V(6).Infof("speak.ToStream LEAVE\n")

	return &result, nil
}

func (c *SpeakClient) ToFile(ctx context.Context, filename string, text string, options interfaces.SpeakOptions) (*api.SpeakResponse, error) {
	klog.V(6).Infof("speak.ToFile ENTER\n")

	result, err := c.ToStream(ctx, text, options)
	if err != nil {
		klog.V(1).Infof("speak.ToStream failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToFile LEAVE\n")
		return nil, err
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		klog.V(1).Infof("os.Create failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToFile LEAVE\n")
		return nil, err
	}
	defer file.Close()

	_, err = file.Write((*result.Stream).Bytes())
	if err != nil {
		klog.V(1).Infof("io.Copy failed. Err: %v\n", err)
		klog.V(6).Infof("speak.ToFile LEAVE\n")
		return nil, err
	}

	result.Stream = nil
	result.Filename = filename

	klog.V(3).Infof("speak.ToFile Succeeded\n")
	klog.V(6).Infof("speak.ToFile LEAVE\n")

	return result, nil
}
