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
)

type PreRecordedTranscriptionOptions struct {
	Alternatives       int      `json:"alternatives" url:"alternatives,omitempty" `
	AnalyzeSentiment   bool     `json:"analyze_sentiment" url:"analyze_sentiment,omitempty" `
	Callback           string   `json:"callback" url:"callback,omitempty" `
	Dates              bool     `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	DetectEntities     bool     `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool     `json:"detect_language" url:"detect_language,omitempty" `
	DetectTopics       bool     `json:"detect_topics" url:"detect_topics,omitempty" `
	Diarize            bool     `json:"diarize" url:"diarize,omitempty" `
	Diarize_version    string   `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation          bool     `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Keywords           []string `json:"keywords" url:"keywords,omitempty" `
	KeywordBoost       string   `json:"keyword_boost" url:"keyword_boost,omitempty" `
	Language           string   `json:"language" url:"language,omitempty" `
	Measurements       bool     `json:"measurements" url:"measurements,omitempty"`
	Model              string   `json:"model" url:"model,omitempty" `
	Multichannel       bool     `json:"multichannel" url:"multichannel,omitempty" `
	Ner                bool     `json:"ner" url:"ner,omitempty" `
	Numbers            bool     `json:"numbers" url:"numbers,omitempty" `
	Numerals           bool     `json:"numerals" url:"numerals,omitempty" ` // Same as Numbers, old name for same option
	Paragraphs         bool     `json:"paragraphs" url:"paragraphs,omitempty" `
	Profanity_filter   bool     `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate          bool     `json:"punctuate" url:"punctuate,omitempty" `
	Redact             bool     `json:"redact" url:"redact,omitempty" `
	Replace            []string `json:"replace" url:"replace,omitempty" `
	Search             []string `json:"search" url:"search,omitempty" `
	Sentiment          bool     `json:"sentiment" url:"sentiment,omitempty" `
	SentimentThreshold float64  `json:"sentiment_threshold" url:"sentiment_threshold,omitempty" `
	Summarize          bool     `json:"summarize" url:"summarize,omitempty" `
	Tag                []string `json:"tag" url:"tag,omitempty"`
	Tier               string   `json:"tier" url:"tier,omitempty" `
	Times              bool     `json:"times" url:"times,omitempty"` // Indicates whether to convert times from written format (e.g., 3:00 pm) to numerical format (e.g., 15:00).
	Translate          string   `json:"translate" url:"translate,omitempty" `
	Utterances         bool     `json:"utterances" url:"utterances,omitempty" `
	Utt_split          int      `json:"utt_split" url:"utt_split,omitempty" `
	Version            string   `json:"version" url:"version,omitempty" `
}

type PreRecordedResponse struct {
	Request_id string   `json:"request_id"`
	Metadata   Metadata `json:"metadata"`
	Results    Results  `json:"results"`
}

type Metadata struct {
	RequestId      string  `json:"request_id"`
	TransactionKey string  `json:"transaction_key"`
	Sha256         string  `json:"sha256"`
	Created        string  `json:"created"`
	Duration       float64 `json:"duration"`
	Channels       int     `json:"channels"`
}

type Hit struct {
	Confidence float64 `json:"confidence"`
	Start      float64 `json:"start"`
	End        float64 `json:"end"`
	Snippet    string  `json:"snippet"`
}

type Search struct {
	Query string `json:"query"`
	Hits  []Hit  `json:"hits"`
}

type WordBase struct {
	Word            string  `json:"word"`
	Start           float64 `json:"start"`
	End             float64 `json:"end"`
	Confidence      float64 `json:"confidence"`
	Punctuated_Word string  `json:"punctuated_word"`
	Speaker         int     `json:"speaker"`
}

type Alternative struct {
	Transcript string     `json:"transcript"`
	Confidence float64    `json:"confidence"`
	Words      []WordBase `json:"words"`
}

type Channel struct {
	Search       []Search      `json:"search"`
	Alternatives []Alternative `json:"alternatives"`
}

type Utterance struct {
	Start      float64    `json:"start"`
	End        float64    `json:"end"`
	Confidence float64    `json:"confidence"`
	Channel    int        `json:"channel"`
	Transcript string     `json:"transcript"`
	Words      []WordBase `json:"words"`
	Speaker    int        `json:"speaker"`
	Id         string     `json:"id"`
}

type Results struct {
	Utterances []Utterance `json:"utterances"`
	Channels   []Channel   `json:"channels"`
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
		return "", errors.New("this function requires a transcript that was generated with the utterances feature")
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
		return "", errors.New("this function requires a transcript that was generated with the utterances feature")
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
