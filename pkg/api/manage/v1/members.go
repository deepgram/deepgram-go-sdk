// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package manage

import (
	"context"
	"log"
	"net/http"

	api "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/manage/v1/interfaces"
	version "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/version"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

func (c *ManageClient) ListMembers(ctx context.Context, projectId string) (*api.MembersResult, error) {
	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersURI, nil, projectId)
	if err != nil {
		// klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}
	log.Printf("Calling %s\n", URI) // TODO

	req, err := http.NewRequestWithContext(ctx, "GET", URI, nil)
	if err != nil {
		// klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MembersResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				// klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				// klog.V(6).Infof("XXXXXXXX LEAVE\n")
				return nil, err
			}
		}

		// klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}

	// klog.V(3).Infof("XXXXXXXX Succeeded\n")
	// klog.V(6).Infof("XXXXXXXX LEAVE\n")
	return &resp, nil
}

func (c *ManageClient) RemoveMember(ctx context.Context, projectId string, memberId string) (*api.MessageResult, error) {
	// checks
	if ctx == nil {
		ctx = context.Background()
	}

	// request
	URI, err := version.GetManageAPI(ctx, c.Client.Options.Host, c.Client.Options.Version, version.MembersByIdURI, nil, projectId, memberId)
	if err != nil {
		// klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}
	log.Printf("Calling %s\n", URI) // TODO

	req, err := http.NewRequestWithContext(ctx, "DELETE", URI, nil)
	if err != nil {
		// klog.V(1).Infof("http.NewRequestWithContext failed. Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}

	// Do it!
	var resp api.MessageResult
	err = c.Client.Do(ctx, req, &resp)

	if err != nil {
		if e, ok := err.(*interfaces.StatusError); ok {
			if e.Resp.StatusCode != http.StatusOK {
				// klog.V(1).Infof("HTTP Code: %v\n", e.Resp.StatusCode)
				// klog.V(6).Infof("XXXXXXXX LEAVE\n")
				return nil, err
			}
		}

		// klog.V(1).Infof("Platform Supplied Err: %v\n", err)
		// klog.V(6).Infof("XXXXXXXX LEAVE\n")
		return nil, err
	}

	// klog.V(3).Infof("XXXXXXXX Succeeded\n")
	// klog.V(6).Infof("XXXXXXXX LEAVE\n")
	return &resp, nil
}
