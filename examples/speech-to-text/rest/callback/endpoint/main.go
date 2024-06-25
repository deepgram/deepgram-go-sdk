// Copyright 2023-2024 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	prettyjson "github.com/hokaccha/go-prettyjson"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/rest/interfaces"
)

type SampleProxy struct {
	server *http.Server
}

func (p *SampleProxy) Start() error {
	// redirect
	router := gin.Default()
	router.POST("/v1/callback", p.postPrerecorded)

	// server
	p.server = &http.Server{
		Addr:    fmt.Sprintf(":3000"),
		Handler: router,
	}

	// start the main entry endpoint to direct traffic
	go func() {
		// this is a blocking call
		err := p.server.ListenAndServeTLS("localhost.crt", "localhost.key")
		if err != nil {
			fmt.Printf("ListenAndServeTLS server stopped. Err: %v\n", err)
		}
	}()

	return nil
}

func (p *SampleProxy) postPrerecorded(c *gin.Context) {
	for key, value := range c.Request.Header {
		fmt.Printf("HTTP Header: %s = %v\n", key, value)
	}

	var prerecordedResponse interfaces.PreRecordedResponse

	// Call BindJSON to bind the received JSON to completionRequest
	if err := c.BindJSON(&prerecordedResponse); err != nil {
		fmt.Printf("BindJSON failed. Err: %v\n", err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bind json failed"})
		return
	}

	data, err := json.Marshal(prerecordedResponse)
	if err != nil {
		fmt.Printf("json.Marshal failed. Err: %v\n", err)
		return
	}

	prettyJSON, err := prettyjson.Format(data)
	if err != nil {
		fmt.Printf("prettyjson.Marshal failed. Err: %v\n", err)
		return
	}

	fmt.Printf("Response:\n%s\n\n", string(prettyJSON))

	c.IndentedJSON(http.StatusOK, "")
}

func (p *SampleProxy) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.server.Shutdown(ctx); err != nil {
		fmt.Printf("Server Shutdown Failed. Err: %v\n", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		fmt.Printf("timeout of 5 seconds.")
	}

	return nil
}

func main() {
	s := &SampleProxy{}
	err := s.Start()
	if err != nil {
		fmt.Printf("SampleProxy.Start failed. Err: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Press ENTER to exit!\n\n")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	err = s.Stop()
	if err != nil {
		fmt.Printf("SampleProxy.Stop failed. Err: %v\n", err)
		os.Exit(1)
	}
}
