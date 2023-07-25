package deepgram

import (
	"encoding/json"
	"net/http"
)

var sdkVersion string = "0.10.0"
var dgAgent string = "deepgram-go-sdk/v" + sdkVersion

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
