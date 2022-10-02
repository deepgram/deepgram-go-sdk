package deepgram

import (
	"encoding/json"
	"net/http"
)


var sdkVersion string = "0.2.2"
var dgAgent string = "deepgram-go-sdk/v" + sdkVersion

type Deepgram struct {
	ApiKey string
	Host string 
	Path string 
}

func Init(apiKey string, host string, path string) *Deepgram {
	if host == "" {
		host = "api.deepgram.com"
	}
	if path == "" {
		path = "/v1/projects"
	}
	return &Deepgram{
		ApiKey: apiKey,
		Host: host,
		Path: path,
	}
}

func GetJson(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(target)
}
