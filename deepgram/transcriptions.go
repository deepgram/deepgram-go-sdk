package deepgram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/gorilla/websocket"
)

type PreRecordedResponse struct {
	Request_id string   `json:"request_id"`
	Metadata   Metadata `json:"metadata"`
	Results    Results  `json:"results"`
}

type LiveTranscriptionOptions struct {
	Alternatives     int      `json:"alternatives" url:"alternatives,omitempty" `
	Callback         string   `json:"callback" url:"callback,omitempty" `
	Channels         int      `json:"channels" url:"channels,omitempty" `
	Dates            bool     `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	Diarize          bool     `json:"diarize" url:"diarize,omitempty" `
	Diarize_version  string   `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation        bool     `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Encoding         string   `json:"encoding" url:"encoding,omitempty" `
	Endpointing      bool     `json:"endpointing" url:"endpointing,omitempty" `
	Interim_results  bool     `json:"interim_results" url:"interim_results,omitempty" `
	Keywords         []string `json:"keywords" url:"keywords,omitempty" `
	Language         string   `json:"language" url:"language,omitempty" `
	Measurements     bool     `json:"measurements" url:"measurements,omitempty" `
	Model            string   `json:"model" url:"model,omitempty" `
	Multichannel     bool     `json:"multichannel" url:"multichannel,omitempty" `
	Ner              bool     `json:"ner" url:"ner,omitempty" `
	Numbers          bool     `json:"numbers" url:"numbers,omitempty" `
	Numerals         bool     `json:"numerals" url:"numerals,omitempty" `
	Profanity_filter bool     `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate        bool     `json:"punctuate" url:"punctuate,omitempty" `
	Redact           bool     `json:"redact" url:"redact,omitempty" `
	Replace          string   `json:"replace" url:"replace,omitempty" `
	Sample_rate      int      `json:"sample_rate" url:"sample_rate,omitempty" `
	Search           []string `json:"search" url:"search,omitempty" `
	Tag              []string `json:"tag" url:"tag,omitempty" `
	Tier             string   `json:"tier" url:"tier,omitempty" `
	Times            bool     `json:"times" url:"times,omitempty" `
	Vad_turnoff      int      `json:"vad_turnoff" url:"vad_turnoff,omitempty" `
	Version          string   `json:"version" url:"version,omitempty" `
}

type PreRecordedTranscriptionOptions struct {
	Alternatives     int      `json:"alternatives" url:"alternatives,omitempty" `
	Callback         string   `json:"callback" url:"callback,omitempty" `
	Dates            bool     `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	DetectLanguage   bool     `json:"detect_language" url:"detect_language,omitempty" `
	DetectTopics     bool     `json:"detect_topics" url:"detect_topics,omitempty" `
	Diarize          bool     `json:"diarize" url:"diarize,omitempty" `
	Diarize_version  string   `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation        bool     `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Keywords         []string `json:"keywords" url:"keywords,omitempty" `
	Language         string   `json:"language" url:"language,omitempty" `
	Measurements     bool     `json:"measurements" url:"measurements,omitempty"`
	Model            string   `json:"model" url:"model,omitempty" `
	Multichannel     bool     `json:"multichannel" url:"multichannel,omitempty" `
	Ner              bool     `json:"ner" url:"ner,omitempty" `
	Numbers          bool     `json:"numbers" url:"numbers,omitempty" `
	Numerals         bool     `json:"numerals" url:"numerals,omitempty" ` // Same as Numbers, old name for same option
	Paragraphs       bool     `json:"paragraphs" url:"paragraphs,omitempty" `
	Profanity_filter bool     `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate        bool     `json:"punctuate" url:"punctuate,omitempty" `
	Redact           bool     `json:"redact" url:"redact,omitempty" `
	Replace          []string `json:"replace" url:"replace,omitempty" `
	Search           []string `json:"search" url:"search,omitempty" `
	Summarize        bool     `json:"summarize" url:"summarize,omitempty" `
	Tag              []string `json:"tag" url:"tag,omitempty"`
	Tier             string   `json:"tier" url:"tier,omitempty" `
	Times            bool     `json:"times" url:"times,omitempty"` // Indicates whether to convert times from written format (e.g., 3:00 pm) to numerical format (e.g., 15:00).
	Utterances       bool     `json:"utterances" url:"utterances,omitempty" `
	Utt_split        int      `json:"utt_split" url:"utt_split,omitempty" `
	Version          string   `json:"version" url:"version,omitempty" `
}

func (dg *Client) LiveTranscription(options LiveTranscriptionOptions) (*websocket.Conn, *http.Response, error) {
	query, _ := query.Values(options)
	u := url.URL{Scheme: "wss", Host: dg.Host, Path: "/v1/listen", RawQuery: query.Encode()}
	log.Printf("connecting to %s", u.String())

	header := http.Header{
		"Host":          []string{dg.Host},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	c, resp, err := websocket.DefaultDialer.Dial(u.String(), header)

	if err != nil {
		log.Printf("handshake failed with status %s", resp.Status)
		log.Fatal("dial:", err)
	}
	return c, resp, nil
}

func (dg *Client) PreRecordedFromURL(source UrlSource, options PreRecordedTranscriptionOptions) (PreRecordedResponse, error) {
	client := new(http.Client)
	query, _ := query.Values(options)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: "/v1/listen", RawQuery: query.Encode()}
	jsonStr, err := json.Marshal(source)
	if err != nil {
		log.Fatal(err)
		return PreRecordedResponse{}, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Fatal(err)
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{"application/json"},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	var result PreRecordedResponse
	res, err := client.Do(req)
	if err != nil {
		return result, err
	}
	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		log.Fatal(string(b))
	}

	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error getting request list: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}

}

func (resp *PreRecordedResponse) ToWebVTT() (string, error) {
	if resp.Results.Utterances == nil {
		return "", errors.New("This function requires a transcript that was generated with the utterances feature.")
	}

	vtt := "WEBVTT\n\n"

	vtt += "NOTE\nTranscription provided by Deepgram\nRequest ID: " + resp.Request_id + "\nCreated: " + resp.Metadata.Created + "\n\n"

	for i, utterance := range resp.Results.Utterances {
		utterance := utterance
		start := SecondsToTimestamp(utterance.Start)
		end := SecondsToTimestamp(utterance.End)
		vtt += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, utterance.Transcript)
	}
	return vtt, nil
}

func (resp *PreRecordedResponse) ToSRT() (string, error) {
	if resp.Results.Utterances == nil {
		return "", errors.New("This function requires a transcript that was generated with the utterances feature.")
	}

	srt := ""

	for i, utterance := range resp.Results.Utterances {
		utterance := utterance
		start := SecondsToTimestamp(utterance.Start)
		end := SecondsToTimestamp(utterance.End)
		end = strings.ReplaceAll(end, ".", ",")
		srt += fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, start, end, utterance.Transcript)

	}
	return srt, nil
}

func SecondsToTimestamp(seconds float64) string {
	hours := int(seconds / 3600)
	minutes := int((seconds - float64(hours*3600)) / 60)
	seconds = seconds - float64(hours*3600) - float64(minutes*60)
	return fmt.Sprintf("%02d:%02d:%02.3f", hours, minutes, seconds)
}
