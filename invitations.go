package deepgram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (dg *Deepgram) ListInvitations(projectId string) (InvitationList, error) {
	client := new(http.Client)
	url := fmt.Sprintf("%s%s/%s/invites", dg.Host(""), dg.Path(""), projectId)

	req , err := http.NewRequest("GET", url, nil)
	if err != nil {
			//Handle Error
			log.Fatal(err)
	}

	req.Header = http.Header{
		"Host": []string{"api.deepgram.com"},
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/" + sdkVersion},
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

func (dg *Deepgram) SendInvitation(projectId string, options InvitationOptions) (Message, error) {
	client := new(http.Client)
	url := fmt.Sprintf("%s%s/%s/invites", dg.Host(""), dg.Path(""), projectId)

	jsonStr, err := json.Marshal(options)
	if err != nil {
		log.Fatal(err)
		return Message{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host": []string{"api.deepgram.com"},
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/" + sdkVersion},
	}

	var result Message
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error sending invitation: %s\n", jsonErr.Error())
		log.Printf("error decoding sakura response: %v", jsonErr)
		if e, ok := err.(*json.SyntaxError); ok {
        log.Printf("syntax error at byte offset %d", e.Offset)
    }
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Deepgram) DeleteInvitation(projectId string, email string) (Message, error) {
	client := new(http.Client)
	url := fmt.Sprintf("%s%s/%s/invites/%s", dg.Host(""), dg.Path(""), projectId, email)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host": []string{"api.deepgram.com"},
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/" + sdkVersion},
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

