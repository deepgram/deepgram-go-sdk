package deepgram

import (
	"encoding/json"
	"net/http"
)

var sdkVersion string = "0.1.2"

type deepgram struct {
	ApiKey string
	Host string 
	Path string 
}

func Init(apiKey string, host string, path string) *deepgram {
	if host == "" {
		host = "api.deepgram.com"
	}
	if path == "" {
		path = "/v1/projects"
	}
	return &deepgram{
		ApiKey: apiKey,
		Host: host,
		Path: path,
	}
}

func GetJson(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()
	
	return json.NewDecoder(resp.Body).Decode(target)
}
