package deepgram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/websocket"
)

type PreRecordedResponse struct {
	Request_id string `json:"request_id"`
	Metadata interface{} `json:"metadata"`
	Results interface{} `json:"results"`
}

func (dg *deepgram) LiveTranscription(options LiveTranscriptionOptions) (*websocket.Conn, *http.Response, error) {
query, _ := query.Values(options)
u := url.URL{Scheme: "wss", Host: dg.Host, Path: "/v1/listen", RawQuery: query.Encode()}
log.Printf("connecting to %s", u.String())

header := http.Header{
		"Host": []string{dg.Host},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/" + sdkVersion},
	}

c, resp, err := websocket.DefaultDialer.Dial(u.String(), header);

if err != nil {
	log.Printf("handshake failed with status %s", resp.Status)
	log.Fatal("dial:", err)
}
return c, resp, nil
  
}

func(dg *deepgram) PreRecordedFromURL(source UrlSource, options PreRecordedTranscriptionOptions) (PreRecordedResponse, error) {
	client := new(http.Client)
	query, _ := query.Values(options)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: "/v1/listen", RawQuery: query.Encode()}
	jsonStr, err := json.Marshal(source)
	if err != nil {
		log.Fatal(err)
		return PreRecordedResponse{}, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
			//Handle Error
			log.Fatal(err)
	}

	req.Header = http.Header{
		"Host": []string{dg.Host},
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/" + sdkVersion},
	}

	var result PreRecordedResponse
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
		fmt.Printf("error getting request list: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
	
}