package interfaces

import (
	"net/http"
	"net/url"
)

// WSClientOptions is the minimal interface required by the WebSocket client.
// Both v1 and v2 ClientOptions implement this interface implicitly.
type ClientOptions interface {
	Parse() error
	SetAPIKey(string)
	GetAuthToken() (string, bool)
	GetHost() string
	GetWSHeaderProcessor() func(http.Header)
	GetRedirectService() bool
	GetSkipServerAuth() bool
	GetProxy() func(*http.Request) (*url.URL, error)
}
