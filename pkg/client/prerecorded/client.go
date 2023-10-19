package prerecorded

type Client struct {
	ApiKey            string
	Host              string
	Path              string
	TranscriptionPath string
}

func New(apiKey string) *Client {
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
