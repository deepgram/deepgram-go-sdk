package deepgram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type UsageRequestListOptions struct {
  Start string `json:"start" url:"start,omitempty"`
  End string `json:"end" url:"end,omitempty"`
  Page int `json:"page" url:"page,omitempty"`
  Limit int `json:"limit" url:"limit,omitempty"`
  Status string `json:"status" url:"status,omitempty"`
};

type UsageRequestList struct {
	Page int `json:"page" url:"page,omitempty"`
	Limit int `json:"limit" url:"limit,omitempty"`
	Requests interface{} `json:"requests" url:"requests,omitempty"`
};
type UsageRequest struct {
	RequestId string `json:"request_id" url:"request_id,omitempty"`
	Created string `json:"created" url:"created,omitempty"`
	Path string `json:"path" url:"path,omitempty"`
	Accessor string `json:"accessor" url:"accessor,omitempty"`
	Response interface{} `json:"response" url:"response,omitempty"`
	Callback interface{} `json:"callback" url:"callback,omitempty"`
};

func (dg *deepgram) ListRequests(projectId string, options UsageRequestListOptions) (UsageRequestList, error) {
	query, _ := query.Values(options)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/requests", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path, RawQuery: query.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
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