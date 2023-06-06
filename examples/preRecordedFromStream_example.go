package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/deepgram-devs/go-sdk/deepgram"
)

var (
	key      = flag.String("key", "", "Deepgram API key")
	file     = flag.String("file", "", "Path to file that will be transcribed")
	mimetype = flag.String("mimetype", "", "Mimetype of file")
)

func main() {
	flag.Parse()

	dg := deepgram.NewClient(*key)
	file, err := os.Open(*file)
	if err != nil {
		log.Fatalf("error opening file %s: %v", *file, err)
	}
	source := deepgram.ReadStreamSource{Stream: file, Mimetype: *mimetype}
	res, err := dg.PreRecordedFromStream(source, deepgram.PreRecordedTranscriptionOptions{Punctuate: true, Diarize: true, Language: "en-US", Utterances: true})
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
