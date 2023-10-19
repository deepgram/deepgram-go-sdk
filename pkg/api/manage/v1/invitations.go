// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package manage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func (dg *ManageClient) ListInvitations(projectId string) (InvitationList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/invites", dg.Client.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result InvitationList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error getting invitation list: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) SendInvitation(projectId string, options InvitationOptions) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/invites", dg.Client.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}
	jsonStr, err := json.Marshal(options)
	if err != nil {
		log.Panic(err)
		return Message{}, err
	}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error sending invitation: %s\n", jsonErr.Error())
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) DeleteInvitation(projectId string, email string) (Message, error) {
	client := new(http.Client)
	// url := fmt.Sprintf("%s%s/%s/invites/%s", dg.Client.Host, dg.Client.Path, projectId, email)
	path := fmt.Sprintf("%s/%s/invites/%s", dg.Client.Path, projectId, email)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Panic(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Client.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"User-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error deleting invitation: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}
