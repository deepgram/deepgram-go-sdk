package main

import (
	"fmt"
	"log"
	"os"

	"github.com/deepgram-devs/go-sdk/deepgram"
)

func main() {
	dg := deepgram.NewClient("YOUR_API_KEY")
	// Feel free to use this as your audio url to test https://anchor.fm/s/3e9db190/podcast/play/22624519/https%3A%2F%2Fd3ctxlq1ktw2nl.cloudfront.net%2Fstaging%2F2020-10-15%2F128822202-44100-1-79cab5de0d7af3c9.mp3
	res, err := dg.PreRecordedFromURL(deepgram.UrlSource{Url: "AUDIO_FILE_URL"}, deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Diarize: true, Language: "en-US", Utterances: true})
	if err != nil {
		fmt.Println("ERROR", err)
		return
	}
	// Log the results
	log.Printf("recv: %+v", res.Results)
	f, err := os.Create("transcription.vtt")
	if err != nil {
		fmt.Printf("error creating VTT file: %v", err)
	}
	// Convert the results to WebVTT format
	vtt, err := res.ToWebVTT()
	if err != nil {
		log.Fatal(err)
	}
	f.WriteString(vtt)

	// Convert the results to SRT format
	srtF, err := os.Create("transcription.srt")
	if err != nil {
		fmt.Printf("error creating SRT file: %v", err)
	}
	srt, err := res.ToSRT()
	if err != nil {
		log.Fatal(err)
	}
	srtF.WriteString(srt)

}
