package deepgram

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"
)

var sdkVersion string = "0.10.0"

var dgAgent string = "@deepgram/sdk/" + sdkVersion + " go/" + goVersion()

type Client struct {
	ApiKey            string
	Host              string
	Path              string
	TranscriptionPath string
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiKey:            apiKey,
		Host:              "api.deepgram.com",
		Path:              "/v1/projects",
		TranscriptionPath: "/v1/listen",
	}
}

func (c *Client) WithHost(host string) *Client {
	c.Host = host
	return c
}

func (c *Client) WithPath(path string) *Client {
	c.Path = path
	return c
}

func (c *Client) WithTranscriptionPath(path string) *Client {
	c.TranscriptionPath = path
	return c
}

func GetJson(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}

func goVersion() string {
	version := runtime.Version()
	if strings.HasPrefix(version, "go") {
		return version[2:]
	}
	return version
}
