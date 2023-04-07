package deepgram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Key struct {
	ApiKeyId string   `json:"api_key_id"`
	Key      string   `json:"key"`
	Comment  string   `json:"comment"`
	Created  string   `json:"created"`
	Scopes   []string `json:"scopes"`
}

type KeyResponseObj struct {
	Member Member `json:"member"`
	ApiKey Key    `json:"api_key"`
}

type KeyResponse struct {
	ApiKeys []KeyResponseObj `json:"api_keys"`
}

type CreateKeyOptions struct {
	ExpirationDate time.Time `json:"expiration_date"`
	TimeToLive     int       `json:"time_to_live"`
	Tags           []string  `json:"tags"`
}

type CreateKeyRequest struct {
	Comment        string   `json:"comment"`
	Scopes         []string `json:"scopes"`
	ExpirationDate string   `json:"expiration_date,omitempty"`
	TimeToLive     int      `json:"time_to_live,omitempty"`
}

func (dg *Client) ListKeys(projectId string) (KeyResponse, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/keys", dg.Path, projectId)
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

	var result KeyResponse
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
		fmt.Printf("error getting keys: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetKey(projectId string, keyId string) (KeyResponseObj, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/keys/%s", dg.Path, projectId, keyId)
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

	var result KeyResponseObj
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
		fmt.Printf("error getting key: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) CreateKey(projectId string, comment string, scopes []string, options CreateKeyOptions) (Key, error) {
	var expirationDate string
	if options.ExpirationDate.IsZero() {
		expirationDate = ""
	} else {
		expirationDate = options.ExpirationDate.Format(time.RFC3339)
	}
	out, err := json.Marshal(CreateKeyRequest{
		Comment:        comment,
		Scopes:         scopes,
		ExpirationDate: expirationDate,
		TimeToLive:     options.TimeToLive,
	})
	fmt.Println(string(out))
	buf := bytes.NewBuffer(out)
	if err != nil {
		log.Fatal(err)
	}
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/keys", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("POST", u.String(), buf)
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

	var result Key
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
		fmt.Printf("error Creating key: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) DeleteKey(projectId string, keyId string) (Message, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/keys/%s", dg.Path, projectId, keyId)
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
		fmt.Printf("error Creating key: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return Message{
			Message: "Key Deleted",
		}, nil
	}
}
