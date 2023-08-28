package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/Jeffail/gabs/v2"
	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
	"github.com/gorilla/websocket"
)

const (
	DEEPGRAM_API_KEY       = "DEEPGRAM_API_KEY"
	STREAM_URL             = "http://stream.live.vc.bbcmedia.co.uk/bbc_world_service"
	CHUNK_SIZE             = 1024 * 2
	TEN_MILLISECONDS_SLEEP = 10 * time.Millisecond
)

func main() {
	client := new(http.Client)

	dg := *deepgram.NewClient(DEEPGRAM_API_KEY)

	res, err := client.Get(STREAM_URL)
	if err != nil {
		log.Println("ERROR getting stream", err)
		return
	}
	defer res.Body.Close()

	fmt.Println("Stream is up and running ", reflect.TypeOf(res))

	reader := bufio.NewReader(res.Body)

	liveTranscriptionOptions := deepgram.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	dgConn, _, err := dg.LiveTranscription(liveTranscriptionOptions)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}
	defer dgConn.Close()

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
