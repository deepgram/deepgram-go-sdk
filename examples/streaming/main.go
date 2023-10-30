package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/dvonthenen/websocket"

	client "github.com/deepgram-devs/deepgram-go-sdk/pkg/client/live"
)

const (
	STREAM_URL             = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"
	CHUNK_SIZE             = 1024 * 2
	TEN_MILLISECONDS_SLEEP = 10 * time.Millisecond
)

func main() {
	var deepgramApiKey string
	if v := os.Getenv("DEEPGRAM_API_KEY"); v != "" {
		log.Println("DEEPGRAM_API_KEY found")
		deepgramApiKey = v
	} else {
		log.Fatal("DEEPGRAM_API_KEY not found")
		os.Exit(1)
	}

	// HTTP client
	httpClient := new(http.Client)

	res, err := httpClient.Get(STREAM_URL)
	if err != nil {
		log.Println("ERROR getting stream", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Stream is up and running ", reflect.TypeOf(res))

	reader := bufio.NewReader(res.Body)

	// live transcription
	liveTranscriptionOptions := client.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	dgConn, _, err := client.New(deepgramApiKey, liveTranscriptionOptions)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}
	defer dgConn.Close()

	// process messages
	chunk := make([]byte, CHUNK_SIZE)
	go func() {
		for {
			_, message, err := dgConn.ReadMessage()
			if err != nil {
				log.Println("ERROR reading message:", err)
				return
			}

			jsonParsed, jsonErr := gabs.ParseJSON(message)
			if jsonErr != nil {
				log.Println("ERROR parsing JSON message:", err)
				return
			}
			log.Printf("recv: %s", jsonParsed.Path("channel.alternatives.0.transcript").String())
		}
	}()

	for {
		bytesRead, err := reader.Read(chunk)

		if err != nil {
			log.Println("ERROR reading chunk:", err)
			return
		}
		err = dgConn.WriteMessage(websocket.BinaryMessage, chunk[:bytesRead])
		if err != nil {
			log.Println("ERROR writing message:", err)
			return
		}
		time.Sleep(TEN_MILLISECONDS_SLEEP)
	}
}
