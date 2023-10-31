// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package prerecorded

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	version "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
	rest "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/rest"
)

type urlSource struct {
	Url string `json:"url"`
}

// New allocated a REST client
func New(apiKey string) *Client {
	if apiKey == "" {
		if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
			log.Println("DEEPGRAM_API_KEY found")
			apiKey = v
		} else {
			log.Println("DEEPGRAM_API_KEY not set")
			return nil
		}
	}

	c := Client{
		Client: rest.New(apiKey),
	}
	return &c
}

// DoFile posts a file capturing a conversation to a given REST endpoint
func (c *Client) DoFile(ctx context.Context, filePath string, req interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	// file?
	fileInfo, err := os.Stat(filePath)
	if err != nil || errors.Is(err, os.ErrNotExist) {
		//klog.V(1).Infof("File %s does not exist. Err : %v\n", filePath, err)
		return err
	}

	if fileInfo.IsDir() || fileInfo.Size() == 0 {
		//klog.V(1).Infof("%s is a directory not a file\n", filePath)
		return ErrInvalidInput
	}

	file, err := os.Open(filePath)
	if err != nil {
		//klog.V(1).Infof("os.Open(%s) failed. Err : %v\n", filePath, err)
		return err
	}
	defer file.Close()

	return c.DoStream(ctx, file, req, resBody)
}

func (c *Client) DoStream(ctx context.Context, src io.Reader, options interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	//klog.V(6).Infof("rest.doCommonFile ENTER\n")

	// obtain URL
	URI, err := version.GetPrerecordedAPI(ctx, options)
	if err != nil {
		log.Printf("version.GetPrerecordedAPI failed. Err: %v\n", err)
		return err
	}
	// TODO: DO NOT PRINT
	log.Printf("Connecting to %s\n", URI)

	req, err := http.NewRequestWithContext(ctx, "POST", URI, src)
	if err != nil {
		//klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
		return err
	}

	if headers, ok := ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				//klog.V(3).Infof("doCommonFile() Custom Header: %s = %s\n", k, v)
				req.Header.Add(k, v)
			}
		}
	}

	req.Header.Set("Host", options.Host)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "token "+c.ApiKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	err = c.HttpClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			//klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
			// detail, errBody := io.ReadAll(res.Body)
			detail, _ := io.ReadAll(res.Body)
			if err != nil {
				//klog.V(4).Infof("io.ReadAll failed. Err: %e\n", errBody)
				//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
				return &interfaces.StatusError{res}
			}
			//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			//klog.V(1).Infof("resBody == nil\n")
			//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			//klog.V(3).Infof("RawResponse\n")
			//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
			return res.Write(b)
		case io.Writer:
			//klog.V(3).Infof("io.Writer\n")
			//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			//klog.V(3).Infof("json.NewDecoder\n")
			d := json.NewDecoder(res.Body)
			//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
			return d.Decode(resBody)
		}
	})

	if err != nil {
		//klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
		return err
	}

	//klog.V(3).Infof("rest.doCommonFile Succeeded\n")
	//klog.V(6).Infof("rest.doCommonFile LEAVE\n")
	return nil
}

// IsUrl returns true if a string is of a URL format
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// DoURL performs a REST call using a URL conversation source
func (c *Client) DoURL(ctx context.Context, url string, options interfaces.PreRecordedTranscriptionOptions, resBody interface{}) error {
	//klog.V(6).Infof("rest.DoURL ENTER\n")
	//klog.V(4).Infof("rest.doCommonURL apiURI: %s\n", apiURI)

	// checks
	validURL := IsUrl(url)
	if !validURL {
		//klog.V(1).Infof("Invalid URL: %s\n", ufRequest.URL)
		//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
		return ErrInvalidInput
	}

	// obtain URL
	URI, err := version.GetPrerecordedAPI(ctx, options)
	if err != nil {
		log.Printf("version.GetPrerecordedAPI failed. Err: %v\n", err)
		return err
	}
	// TODO: DO NOT PRINT
	log.Printf("Connecting to %s\n", URI)

	var buf bytes.Buffer
	err = json.NewEncoder(&buf).Encode(urlSource{Url: url})
	if err != nil {
		//klog.V(1).Infof("json.NewEncoder().Encode() failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", URI, &buf)
	if err != nil {
		//klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
		return err
	}

	if headers, ok := ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				//klog.V(3).Infof("doCommonURL() Custom Header: %s = %s\n", k, v)
				req.Header.Add(k, v)
			}
		}
	}

	req.Header.Set("Host", options.Host)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.ApiKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		//klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	err = c.HttpClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			//klog.V(4).Infof("HTTP Error Code: %d\n", res.StatusCode)
			// detail, errBody := io.ReadAll(res.Body)
			detail, _ := io.ReadAll(res.Body)
			if err != nil {
				//klog.V(1).Infof("io.ReadAll failed. Err: %e\n", errBody)
				//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
				return &interfaces.StatusError{res}
			}
			//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			//klog.V(1).Infof("resBody == nil\n")
			//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			//klog.V(3).Infof("RawResponse\n")
			//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
			return res.Write(b)
		case io.Writer:
			//klog.V(3).Infof("io.Writer\n")
			//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			//klog.V(3).Infof("json.NewDecoder\n")
			d := json.NewDecoder(res.Body)
			//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
			return d.Decode(resBody)
		}
	})

	if err != nil {
		//klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
		return err
	}

	//klog.V(3).Infof("rest.doCommonURL Succeeded\n")
	//klog.V(6).Infof("rest.doCommonURL LEAVE\n")
	return nil
}

// Do is a generic REST API call to the platform
func (c *Client) Do(ctx context.Context, req *http.Request, resBody interface{}) error {
	//klog.V(6).Infof("rest.Do ENTER\n")

	if headers, ok := ctx.Value(interfaces.HeadersContext{}).(http.Header); ok {
		for k, v := range headers {
			for _, v := range v {
				//klog.V(3).Infof("Do() Custom Header: %s = %s\n", k, v)
				req.Header.Add(k, v)
			}
		}
	}

	// req.Header.Set("Host", c.options.Host)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "token "+c.ApiKey)
	req.Header.Set("User-Agent", interfaces.DgAgent)

	switch req.Method {
	case http.MethodPost, http.MethodPatch, http.MethodPut:
		//klog.V(3).Infof("Content-Type = application/json\n")
		req.Header.Set("Content-Type", "application/json")
	}

	err := c.HttpClient.Do(ctx, req, func(res *http.Response) error {
		switch res.StatusCode {
		case http.StatusOK:
		case http.StatusCreated:
		case http.StatusNoContent:
		case http.StatusBadRequest:
			//klog.V(1).Infof("HTTP Error Code: %d\n", res.StatusCode)
			detail, errBody := io.ReadAll(res.Body)
			if errBody != nil {
				//klog.V(1).Infof("io.ReadAll failed. Err: %e\n", errBody)
				//klog.V(6).Infof("rest.DoFile LEAVE\n")
				return &interfaces.StatusError{res}
			}
			//klog.V(6).Infof("rest.Do LEAVE\n")
			return fmt.Errorf("%s: %s", res.Status, bytes.TrimSpace(detail))
		default:
			return &interfaces.StatusError{res}
		}

		if resBody == nil {
			//klog.V(1).Infof("resBody == nil\n")
			//klog.V(6).Infof("rest.Do LEAVE\n")
			return nil
		}

		switch b := resBody.(type) {
		case *interfaces.RawResponse:
			//klog.V(3).Infof("RawResponse\n")
			//klog.V(6).Infof("rest.Do LEAVE\n")
			return res.Write(b)
		case io.Writer:
			//klog.V(3).Infof("io.Writer\n")
			//klog.V(6).Infof("rest.Do LEAVE\n")
			_, err := io.Copy(b, res.Body)
			return err
		default:
			//klog.V(3).Infof("json.NewDecoder\n")
			d := json.NewDecoder(res.Body)
			//klog.V(6).Infof("rest.Do LEAVE\n")
			return d.Decode(resBody)
		}
	})

	if err != nil {
		//klog.V(1).Infof("err = c.Client.Do failed. Err: %v\n", err)
		//klog.V(6).Infof("rest.Do LEAVE\n")
		return err
	}

	//klog.V(3).Infof("rest.Do Succeeded\n")
	//klog.V(6).Infof("rest.Do LEAVE\n")
	return nil
}
