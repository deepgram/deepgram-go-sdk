// Copyright 2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package daily_test

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	prerecorded "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen/v1/rest"

	utils "github.com/deepgram/deepgram-go-sdk/tests/utils"
)

const (
	model string = "2-general-nova"

	url string = "https://static.deepgram.com/examples/Bueller-Life-moves-pretty-fast.wav"
)

const (
	FromURLSmartFormat = "Yep. I said it before, and I'll say it again. Life moves pretty fast. You don't stop and look around once in a while, you could miss it."
	FromURLSummarize   = "Yep. I said it before, and I'll say it again. Life moves pretty fast. You don't stop and look around once in a while, you could miss it."
)

func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..") // change to suit test file location
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestDaily_Prerecorded(t *testing.T) {
	t.Run("FromURL-SmartFormat", func(t *testing.T) {
		// context
		ctx := context.Background()

		// send stream to Deepgram
		options := &interfaces.PreRecordedTranscriptionOptions{
			Model:       "nova-2",
			SmartFormat: true,
		}

		// create a Deepgram client
		c := client.NewWithDefaults()
		dg := prerecorded.New(c)

		// send the URL to Deepgram
		res, err := dg.FromURL(ctx, url, options)
		if err != nil {
			t.Errorf("FromURL(%v) Err: %v", options, err)
		}

		// save the options
		data, err := json.Marshal(options)
		if err != nil {
			t.Errorf("json.Marshal Err: %v", err)
		}

		sha256sum := fmt.Sprintf("%x", sha256.Sum256(data))
		filename := fmt.Sprintf("tests/response_data/%s.cmd", strings.ReplaceAll(t.Name(), "/", "_"))
		err = utils.SaveMetadataString(filename, sha256sum)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		filename = fmt.Sprintf("tests/response_data/%s-options.json", sha256sum)
		err = utils.SaveMetadataBytes(filename, data)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		// save the response
		data, err = json.Marshal(res)
		if err != nil {
			t.Errorf("json.Marshal Err: %v", err)
		}

		filename = fmt.Sprintf("tests/response_data/%s-response.json", sha256sum)
		err = utils.SaveMetadataBytes(filename, data)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		// check the response
		for _, value := range res.Metadata.ModelInfo {
			if strings.Compare(model, value.Name) != 0 {
				t.Errorf("%s: %s != %s", t.Name(), model, value.Name)
			}
		}
		transcript := res.Results.Channels[0].Alternatives[0].Transcript
		if strings.Compare(FromURLSmartFormat, transcript) != 0 {
			t.Errorf("%s: %s != %s", t.Name(), FromURLSmartFormat, transcript)
		}
	})
}

func TestDaily_Prerecorded_summary(t *testing.T) {
	t.Run("FromURL-Summarize", func(t *testing.T) {
		// context
		ctx := context.Background()

		// send stream to Deepgram
		options := &interfaces.PreRecordedTranscriptionOptions{
			Model:     "nova-2",
			Summarize: "v2",
		}

		// create a Deepgram client
		c := client.NewWithDefaults()
		dg := prerecorded.New(c)

		// send the URL to Deepgram
		res, err := dg.FromURL(ctx, url, options)
		if err != nil {
			t.Errorf("FromURL(%v) Err: %v", options, err)
		}

		// save the options
		data, err := json.Marshal(options)
		if err != nil {
			t.Errorf("json.Marshal Err: %v", err)
		}

		sha256sum := fmt.Sprintf("%x", sha256.Sum256(data))
		filename := fmt.Sprintf("tests/response_data/%s.cmd", strings.ReplaceAll(t.Name(), "/", "_"))
		err = utils.SaveMetadataString(filename, sha256sum)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		filename = fmt.Sprintf("tests/response_data/%s-options.json", sha256sum)
		err = utils.SaveMetadataBytes(filename, data)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		// save the response
		data, err = json.Marshal(res)
		if err != nil {
			t.Errorf("json.Marshal Err: %v", err)
		}

		filename = fmt.Sprintf("tests/response_data/%s-response.json", sha256sum)
		err = utils.SaveMetadataBytes(filename, data)
		if err != nil {
			t.Errorf("dataFile.WriteString Err: %v", err)
		}

		// check the response
		for _, value := range res.Metadata.ModelInfo {
			if strings.Compare(model, value.Name) != 0 {
				t.Errorf("%s: %s != %s", t.Name(), model, value.Name)
			}
		}
		summary := res.Results.Summary.Short
		if strings.Compare(FromURLSummarize, summary) != 0 {
			t.Errorf("%s: %s != %s", t.Name(), FromURLSummarize, summary)
		}
	})
}
