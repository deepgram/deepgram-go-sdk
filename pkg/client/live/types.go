package live

import (
	"context"
	"sync"

	"github.com/dvonthenen/websocket"

	live "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1"
	msginterface "github.com/deepgram-devs/deepgram-go-sdk/pkg/api/live/v1/interfaces"
	interfaces "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/interfaces"
)

// Credentials for connecting to the platform
type Credentials struct {
	Host            string
	ApiKey          string
	RedirectService bool
	SkipServerAuth  bool
}

// Client return websocket client connection
type Client struct {
	options interfaces.LiveTranscriptionOptions

	sendBuf   chan []byte
	org       context.Context
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.RWMutex
	wsconn *websocket.Conn
	retry  bool

	creds    *Credentials
	callback msginterface.LiveMessageCallback
	router   *live.MessageRouter
}
