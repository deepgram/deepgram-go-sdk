package deepgram

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

func (dg *Deepgram) LiveTranscription(options LiveTranscriptionOptions) (*websocket.Conn, *http.Response, error) {

u := url.URL{Scheme: "wss", Host: "api.deepgram.com", Path: "/v1/listen",}
log.Printf("connecting to %s", u.String())

header := http.Header{
		"Host": []string{"api.deepgram.com"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent": []string{"go-sdk/1.0.0"},
	}

c, resp, err := websocket.DefaultDialer.Dial(u.String(), header);

if err != nil {
	log.Printf("handshake failed with status %d", resp.StatusCode)
	log.Fatal("dial:", err)
}
return c, resp, nil
  
}

