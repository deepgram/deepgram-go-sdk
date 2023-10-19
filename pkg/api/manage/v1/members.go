// Copyright 2023 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package manage

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Member struct {
	Email     string   `json:"email"`
	MemberId  string   `json:"member_id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Scopes    []string `json:"scopes"`
}

type MemberList struct {
	Members []Member `json:"members"`
}

type ScopeList struct {
	Scopes []string `json:"scopes"`
}

func (dg *ManageClient) ListMembers(projectId string) (MemberList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members", dg.Client.Path, projectId)
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

	var result MemberList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting members: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) RemoveMember(projectId string, memberId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s", dg.Client.Path, projectId, memberId)
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
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error removing member: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) GetMemberScopes(projectId string, memberId string) (ScopeList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s/scopes", dg.Client.Path, projectId, memberId)
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

	var result ScopeList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting member scopes: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

type MemberScope struct {
	Scope string `json:"scope"`
}

func (dg *ManageClient) UpdateMemberScopes(projectId string, memberId string, scope string) (Message, error) {
	newScope := fmt.Sprintf(`{"scope":"%s"}`, scope)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s/scopes", dg.Client.Path, projectId, memberId)
	u := url.URL{Scheme: "https", Host: dg.Client.Host, Path: path}
	req, err := http.NewRequest("PUT", u.String(), strings.NewReader(newScope))
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
	if res.StatusCode != 200 {

		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error updating member scopes: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *ManageClient) LeaveProject(projectId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/leave", dg.Client.Path, projectId)
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
	if res.StatusCode != 200 {

		b, _ := io.ReadAll(res.Body)
		log.Panic(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error leaving project: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}
