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
	Alternatives       int         `json:"alternatives" url:"alternatives,omitempty" `
	AnalyzeSentiment   bool        `json:"analyze_sentiment" url:"analyze_sentiment,omitempty" `
	Callback           string      `json:"callback" url:"callback,omitempty" `
	Dates              bool        `json:"dates" url:"dates,omitempty"` // Indicates whether to convert dates from written format (e.g., january first) to numerical format (e.g., 01-01).
	DetectEntities     bool        `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool        `json:"detect_language" url:"detect_language,omitempty" `
	DetectTopics       bool        `json:"detect_topics" url:"detect_topics,omitempty" `
	Diarize            bool        `json:"diarize" url:"diarize,omitempty" `
	Diarize_version    string      `json:"diarize_version" url:"diarize_version,omitempty" `
	Dictation          bool        `json:"dictation" url:"dictation,omitempty"` // Option to format punctuated commands. Eg: "i went to the store period new paragraph then i went home" --> "i went to the store. <\n> then i went home"
	Keywords           []string    `json:"keywords" url:"keywords,omitempty" `
	KeywordBoost       string      `json:"keyword_boost" url:"keyword_boost,omitempty" `
	Language           string      `json:"language" url:"language,omitempty" `
	Measurements       bool        `json:"measurements" url:"measurements,omitempty"`
	Model              string      `json:"model" url:"model,omitempty" `
	Multichannel       bool        `json:"multichannel" url:"multichannel,omitempty" `
	Ner                bool        `json:"ner" url:"ner,omitempty" `
	Numbers            bool        `json:"numbers" url:"numbers,omitempty" `
	Numerals           bool        `json:"numerals" url:"numerals,omitempty" ` // Same as Numbers, old name for same option
	Paragraphs         bool        `json:"paragraphs" url:"paragraphs,omitempty" `
	Profanity_filter   bool        `json:"profanity_filter" url:"profanity_filter,omitempty" `
	Punctuate          bool        `json:"punctuate" url:"punctuate,omitempty" `
	Redact             bool        `json:"redact" url:"redact,omitempty" `
	Replace            []string    `json:"replace" url:"replace,omitempty" `
	Search             []string    `json:"search" url:"search,omitempty" `
	Sentiment          bool        `json:"sentiment" url:"sentiment,omitempty" `
	SentimentThreshold float64     `json:"sentiment_threshold" url:"sentiment_threshold,omitempty" `
	SmartFormat        bool        `json:"smart_format" url:"smart_format,omitempty" `
	Summarize          interface{} `json:"summarize" url:"summarize,omitempty" ` // bool | string
	Tag                []string    `json:"tag" url:"tag,omitempty"`
	Tier               string      `json:"tier" url:"tier,omitempty" `
	Times              bool        `json:"times" url:"times,omitempty"` // Indicates whether to convert times from written format (e.g., 3:00 pm) to numerical format (e.g., 15:00).
	Translate          string      `json:"translate" url:"translate,omitempty" `
	Utterances         bool        `json:"utterances" url:"utterances,omitempty" `
	Utt_split          float64     `json:"utt_split" url:"utt_split,omitempty" `
	Version            string      `json:"version" url:"version,omitempty" `
	FillerWords      string   `json:"filler_words" url:"filler_words,omitempty" `
}

type PreRecordedResponse struct {
	RequestId string   `json:"request_id"`
	Metadata  Metadata `json:"metadata"`
	Results   Results  `json:"results"`
}

type Metadata struct {
	TransactionKey string   `json:"transaction_key"`
	RequestId      string   `json:"request_id"`
	Sha256         string   `json:"sha256"`
	Created        string   `json:"created"`
	Duration       float64  `json:"duration"`
	Channels       int      `json:"channels"`
	Models         []string `json:"models"`
	ModelInfo      map[string]struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Arch    string `json:"arch"`
	} `json:"model_info"`
	Warnings []*Warning `json:"warnings,omitempty"`
}

type Warning struct {
	Parameter string `json:"parameter"`
	Type      string `json:"type"`
	Message   string `json:"message"`
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
	Word              string  `json:"word"`
	Start             float64 `json:"start"`
	End               float64 `json:"end"`
	Confidence        float64 `json:"confidence"`
	Speaker           *int    `json:"speaker,omitempty"`
	SpeakerConfidence float64 `json:"speaker_confidence,omitempty"`
	Punctuated_Word   string  `json:"punctuated_word,omitempty"`
	Sentiment         string  `json:"sentiment,omitempty"`
}

type Alternative struct {
	Transcript string          `json:"transcript"`
	Confidence float64         `json:"confidence"`
	Words      []WordBase      `json:"words"`
	Summaries  []*SummaryV1    `json:"summaries,omitempty"`
	Paragraphs *ParagraphGroup `json:"paragraphs,omitempty"`
	Topics     []*TopicBase    `json:"topics,omitempty"`
	Entities   []*EntityBase   `json:"entities,omitempty"`
}

type ParagraphGroup struct {
	Transcript string          `json:"transcript"`
	Paragraphs []ParagraphBase `json:"paragraphs"`
}

type ParagraphBase struct {
	Sentences []SentenceBase `json:"sentences"`
	NumWords  int            `json:"num_words"`
	Start     float64        `json:"start"`
	End       float64        `json:"end"`
}

type SentenceBase struct {
	Text  string  `json:"text"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

type EntityBase struct {
	Label      string  `json:"label"`
	Value      string  `json:"value"`
	Confidence float64 `json:"confidence"`
	StartWord  int     `json:"start_word"`
	EndWord    int     `json:"end_word"`
}

type TopicBase struct {
	Text      string  `json:"text"`
	StartWord int     `json:"start_word"`
	EndWord   int     `json:"end_word"`
	Topics    []Topic `json:"topics"`
}

type Topic struct {
	Topic      string  `json:"topic"`
	Confidence float64 `json:"confidence"`
}

type Channel struct {
	Search           []*Search     `json:"search,omitempty"`
	Alternatives     []Alternative `json:"alternatives"`
	DetectedLanguage string        `json:"detected_language,omitempty"`
}

type Utterance struct {
	Start      float64    `json:"start"`
	End        float64    `json:"end"`
	Confidence float64    `json:"confidence"`
	Channel    int        `json:"channel"`
	Transcript string     `json:"transcript"`
	Words      []WordBase `json:"words"`
	Speaker    *int       `json:"speaker,omitempty"`
	Id         string     `json:"id"`
}

type Results struct {
	Utterances []*Utterance `json:"utterances,omitempty"`
	Channels   []Channel    `json:"channels"`
	Summary    *SummaryV2   `json:"summary,omitempty"`
}

type SummaryV1 struct {
	Summary   string `json:"summary"`
	StartWord int    `json:"start_word"`
	EndWord   int    `json:"end_word"`
}

type SummaryV2 struct {
	Short  string `json:"short"`
	Result string `json:"result"`
}

func (dg *Client) PreRecordedFromStream(source ReadStreamSource, options PreRecordedTranscriptionOptions) (*PreRecordedResponse, error) {
	client := &http.Client{}
	query, _ := query.Values(options)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: dg.TranscriptionPath, RawQuery: query.Encode()}

	// TODO: accept file path as string build io.Reader here
	req, err := http.NewRequest("POST", u.String(), source.Stream)
	if err != nil {
		//Handle Error
		return nil, err
	}

	req.Header = http.Header{
		"Host":          []string{dg.Host},
		"Content-Type":  []string{source.Mimetype},
		"Authorization": []string{"token " + dg.ApiKey},
		"X-DG-Agent":    []string{dgAgent},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("response error: %s", string(b))
	}

	var result PreRecordedResponse
	jsonErr := GetJson(res, &result)
	if jsonErr != nil {
		fmt.Printf("error getting request list: %s\n", jsonErr.Error())
		return nil, jsonErr
	}

	return &result, nil
}

func (dg *Client) PreRecordedFromURL(source UrlSource, options PreRecordedTranscriptionOptions) (PreRecordedResponse, error) {
	client := new(http.Client)
	query, _ := query.Values(options)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: dg.TranscriptionPath, RawQuery: query.Encode()}
	jsonStr, err := json.Marshal(source)
	if err != nil {
		log.Panic(err)
		return PreRecordedResponse{}, err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonStr))
	if err != nil {
		//Handle Error
		log.Panic(err)
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
		log.Panic(string(b))
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

	vtt += "NOTE\nTranscription provided by Deepgram\nRequest ID: " + resp.Metadata.RequestId + "\nCreated: " + resp.Metadata.Created + "\n\n"

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
