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

func main() {
	client := new(http.Client)
	// IMPORTANT: Make sure you add your own API key here
	dg := *deepgram.NewClient("YOUR_API_KEY")
	resp, err := client.Get("http://stream.live.vc.bbcmedia.co.uk/bbc_world_service")
	if err != nil {
		log.Println("ERRROR getting stream", err)
	}
	fmt.Println("Stream is up and running ", reflect.TypeOf(resp))
	reader := bufio.NewReader(resp.Body)

	options := deepgram.LiveTranscriptionOptions{
		Language:  "en-US",
		Punctuate: true,
	}

	dgConn, _, err := dg.LiveTranscription(options)

	chunk := make([]byte, 1024*2)

	go func() {
		for {
			_, message, err := dgConn.ReadMessage()
			if err != nil {
				fmt.Println("ERROR reading message")
				log.Panic(err)
			}

			jsonParsed, jsonErr := gabs.ParseJSON(message)
			if jsonErr != nil {
				log.Panic(err)
			}
			log.Printf("recv: %s", jsonParsed.Path("channel.alternatives.0.transcript").String())

		}
	}()

	for {
		bytesRead, err := reader.Read(chunk)

		if err != nil {
			fmt.Println("ERROR reading chunk")
			log.Panic(err)
		}
		dgConn.WriteMessage(websocket.BinaryMessage, chunk[:bytesRead])
		time.Sleep(10 * time.Millisecond)

	}

}
