package deepgram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type UsageRequestListOptions struct {
  Start string `json:"start"`
  End string `json:"end"`
  Page int `json:"page"`
  Limit int `json:"limit"`
  Status string `json:"status"`
};

type UsageRequestList struct {
	Page int `json:"page"`
	Limit int `json:"limit"`
	Requests interface{} `json:"requests"`
};
type UsageRequest struct {
	RequestId string `json:"request_id"`
	Created string `json:"created"`
	Path string `json:"path"`
	Accessor string `json:"accessor"`
	Response interface{} `json:"response"`
	Callback interface{} `json:"callback"`
};

func (dg *Deepgram) ListRequests(projectId string, options UsageRequestListOptions) (UsageRequestList, error) {
	client := new(http.Client)
	url := fmt.Sprintf("%s%s/%s/requests", dg.Host(""), dg.Path(""), projectId)
	jsonStr, err := json.Marshal(options)
	if err != nil {
		log.Fatal(err)
		return UsageRequestList{}, err
	}

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	if err != nil {
			//Handle Error
			log.Fatal(err)
	}

	req.Header = http.Header{
		"Host": []string{dg.Host("")},
		"Content-Type": []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/1.0.0"},
	}

	var result UsageRequestList
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