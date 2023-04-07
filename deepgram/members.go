package deepgram

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

func (dg *Client) ListMembers(projectId string) (MemberList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result MemberList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error getting members: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) RemoveMember(projectId string, memberId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s", dg.Path, projectId, memberId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error removing member: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetMemberScopes(projectId string, memberId string) (ScopeList, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s/scopes", dg.Path, projectId, memberId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result ScopeList
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
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

func (dg *Client) UpdateMemberScopes(projectId string, memberId string, scope string) (Message, error) {
	newScope := fmt.Sprintf(`{"scope":"%s"}`, scope)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/members/%s/scopes", dg.Path, projectId, memberId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("PUT", u.String(), strings.NewReader(newScope))
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {

		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error updating member scopes: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) LeaveProject(projectId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/leave", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {

		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}
	jsonErr := GetJson(res, &result)

	if jsonErr != nil {
		fmt.Printf("error leaving project: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}
