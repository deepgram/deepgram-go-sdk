package manage

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strings"

	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/prerecorded"
)

var sdkVersion string = "0.10.0"
var dgAgent string = "@deepgram/sdk/" + sdkVersion + " go/" + goVersion()

func goVersion() string {
	version := runtime.Version()
	if strings.HasPrefix(version, "go") {
		return version[2:]
	}
	return version
}

type ManageClient struct {
	*client.Client
}

func New(client *client.Client) *ManageClient {
	return &ManageClient{client}
}

func GetJson(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(target)
}
