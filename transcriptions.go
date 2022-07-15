package deepgram

import (
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/websocket"
)

func (dg *Deepgram) LiveTranscription(options LiveTranscriptionOptions) (*websocket.Conn, *http.Response, error) {
query, _ := query.Values(options)
u := url.URL{Scheme: "wss", Host: "api.deepgram.com", Path: "/v1/listen", RawQuery: query.Encode()}
log.Printf("connecting to %s", u.String())

header := http.Header{
		"Host": []string{"api.deepgram.com"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/0.1.0"},
	}

c, resp, err := websocket.DefaultDialer.Dial(u.String(), header);

if err != nil {
	log.Printf("handshake failed with status %s", resp.Status)
	log.Fatal("dial:", err)
}
return c, resp, nil
  
}

func(dg *Deepgram) PreRecordedFromURL(options PreRecordedTranscriptionOptions) {
	query, _ := query.Values(options)
	u := url.URL{Scheme: "https", Host: "api.deepgram.com", Path: "/v1/projects", RawQuery: query.Encode()}
	log.Printf("connecting to %s", u.String())

	header := http.Header{
		"Host": []string{"api.deepgram.com"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/0.1.0"},
	}

	
}