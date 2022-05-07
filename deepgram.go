package deepgram

import (
	"encoding/json"
	"net/http"
)

type Deepgram struct {
	ApiKey string
}

func GetJson(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func (dg *Deepgram) Host(host string) string  {
	if host == "" {
		return "https://api.deepgram.com"
	} else {
		return host
	}
}

func (dg *Deepgram) Path(path string) string  {
	if path == "" {
		return "/v1/projects"
	} else {
		return path
	}
}