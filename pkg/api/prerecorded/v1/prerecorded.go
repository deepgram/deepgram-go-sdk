// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package prerecorded

import (
	"context"
	"io"
	"net/http"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/prerecorded/v1/interfaces"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

type PrerecordedClient struct {
	*client.Client
}

func New(client *client.Client) *PrerecordedClient {
	return &PrerecordedClient{client}
}

func (c *PrerecordedClient) FromFile(ctx context.Context, file string, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	// klog.V(6).Infof("FromFile ENTER\n")
	// klog.V(3).Infof("filePath: %s\n", filePath)

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// send the file!
	var resp api.PreRecordedResponse

	err := c.Client.DoFile(ctx, file, options, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				// klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				// klog.V(6).Infof("FromFile LEAVE\n")
				return nil, err
			}
		}

		// klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		// klog.V(6).Infof("FromFile LEAVE\n")
		return nil, err
	}

	// klog.V(3).Infof("FromFile Succeeded\n")
	// klog.V(6).Infof("FromFile LEAVE\n")
	return &resp, nil
}

func (c *PrerecordedClient) FromStream(ctx context.Context, src io.Reader, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	// klog.V(6).Infof("FromStream ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// send the file!
	var resp api.PreRecordedResponse

	err := c.Client.DoStream(ctx, src, options, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				// klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				// klog.V(6).Infof("FromStream LEAVE\n")
				return nil, err
			}
		}

		// klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		// klog.V(6).Infof("FromStream LEAVE\n")
		return nil, err
	}

	// klog.V(3).Infof("FromStream Succeeded\n")
	// klog.V(6).Infof("FromStream LEAVE\n")
	return &resp, nil
}

func (c *PrerecordedClient) FromURL(ctx context.Context, url string, options interfaces.PreRecordedTranscriptionOptions) (*api.PreRecordedResponse, error) {
	// klog.V(6).Infof("FromURL ENTER\n")

	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// send the file!
	var resp api.PreRecordedResponse

	err := c.Client.DoURL(ctx, url, options, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				// klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				// klog.V(6).Infof("FromURL LEAVE\n")
				return nil, err
			}
		}

		// klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		// klog.V(6).Infof("FromURL LEAVE\n")
		return nil, err
	}

	// klog.V(3).Infof("FromURL Succeeded\n")
	// klog.V(6).Infof("FromURL LEAVE\n")

	return &resp, nil
}

// func (resp *PreRecordedResponse) ToWebVTT() (string, error) {
// 	if resp.Results.Utterances == nil {
// 		return "", errors.New("this function requires a transcript that was generated with the utterances feature")
// 	}

// 	vtt := "WEBVTT\n\n"

// 	vtt += "NOTE\nTranscription provided by Deepgram\nRequest ID: " + resp.Metadata.RequestId + "\nCreated: " + resp.Metadata.Created + "\n\n"

// 	for i, utterance := range resp.Results.Utterances {
// 		utterance := utterance
// 		start := SecondsToTimestamp(utterance.Start)
// 		end := SecondsToTimestamp(utterance.End)
// 		vtt += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, utterance.Transcript)
// 	}
// 	return vtt, nil
// }

// func (resp *PreRecordedResponse) ToSRT() (string, error) {
// 	if resp.Results.Utterances == nil {
// 		return "", errors.New("this function requires a transcript that was generated with the utterances feature")
// 	}

// 	srt := ""

// 	for i, utterance := range resp.Results.Utterances {
// 		utterance := utterance
// 		start := SecondsToTimestamp(utterance.Start)
// 		end := SecondsToTimestamp(utterance.End)
// 		end = strings.ReplaceAll(end, ".", ",")
// 		srt += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, utterance.Transcript)

// 	}
// 	return srt, nil
// }

// func SecondsToTimestamp(seconds float64) string {
// 	hours := int(seconds / 3600)
// 	minutes := int((seconds - float64(hours*3600)) / 60)
// 	seconds = seconds - float64(hours*3600) - float64(minutes*60)
// 	return fmt.Sprintf("%02d:%02d:%02.3f", hours, minutes, seconds)
// }

// func GetJson(resp *http.Response, target interface{}) error {
// 	defer resp.Body.Close()

// 	return json.NewDecoder(resp.Body).Decode(target)
// }
