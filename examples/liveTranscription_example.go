package live_example

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	deepgram "github.com/deepgram-devs/go-sdk"
	"github.com/gorilla/websocket"
)

func main() {
	client := new(http.Client)
	// IMPORTANT: Make sure you add your own API key here
	dg := deepgram.Deepgram{
        ApiKey: "YOUR_API_KEY",
    }
	resp, err := client.Get("http://stream.live.vc.bbcmedia.co.uk/bbc_radio_fourlw_online_nonuk")
	if err != nil {
			log.Println("ERRROR getting stream", err)
	}
	fmt.Println("Stream is up and running ", reflect.TypeOf(resp))
	reader := bufio.NewReader(resp.Body)

	options := deepgram.LiveTranscriptionOptions{}
	options.Punctuate = true
	options.Language = "en-US"

	dgConn, _, err := dg.LiveTranscription(options)

	chunk := make([]byte, 1024*2)

	go func() {
	for {
		_, message, err := dgConn.ReadMessage()
		if err != nil {
				fmt.Println("ERROR reading message")
				log.Fatal(err)
		}
		log.Printf("recv: %s", message)
	}
	}()

	for {
		bytesRead, err := reader.Read(chunk)

		if err != nil {
				fmt.Println("ERROR reading chunk")
				log.Fatal(err)
		}
		dgConn.WriteMessage(websocket.BinaryMessage, chunk[:bytesRead])
		time.Sleep(10 * time.Millisecond)

	
	}
	
	defer resp.Body.Close()
	log.Print(resp.Body)
}

