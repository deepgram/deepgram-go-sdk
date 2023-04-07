package deepgram

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

type UsageRequestListOptions struct {
	Start  string `json:"start" url:"start,omitempty"`
	End    string `json:"end" url:"end,omitempty"`
	Page   int    `json:"page" url:"page,omitempty"`
	Limit  int    `json:"limit" url:"limit,omitempty"`
	Status string `json:"status" url:"status,omitempty"`
}

type UsageOptions struct {
	Accessor           string   `json:"accessor" url:"accessor,omitempty"`
	Alternatives       bool     `json:"alternatives" url:"alternatives,omitempty"`
	AnalyzeSentiment   bool     `json:"analyze_sentiment" url:"analyze_sentiment,omitempty"`
	DetectEntities     bool     `json:"detect_entities" url:"detect_entities,omitempty"`
	DetectLanguage     bool     `json:"detect_language" url:"detect_language,omitempty"`
	DetectTopics       bool     `json:"detect_topics" url:"detect_topics,omitempty"`
	Diarize            bool     `json:"diarize" url:"diarize,omitempty"`
	End                string   `json:"end" url:"end,omitempty"`
	InterimResults     bool     `json:"interim_results" url:"interim_results,omitempty"`
	Keywords           bool     `json:"keywords" url:"keywords,omitempty"`
	Method             string   `json:"method" url:"method,omitempty"` // Must be one of "sync" | "async" | "streaming"
	Model              string   `json:"model" url:"model,omitempty"`
	Multichannel       bool     `json:"multichannel" url:"multichannel,omitempty"`
	Ner                bool     `json:"ner" url:"ner,omitempty"`
	Numbers            bool     `json:"numbers" url:"numbers,omitempty"`
	Numerals           bool     `json:"numerals" url:"numerals,omitempty"`
	Paragraphs         bool     `json:"paragraphs" url:"paragraphs,omitempty"`
	ProfanityFilter    bool     `json:"profanity_filter" url:"profanity_filter,omitempty"`
	Punctuate          bool     `json:"punctuate" url:"punctuate,omitempty"`
	Redact             bool     `json:"redact" url:"redact,omitempty"`
	Replace            bool     `json:"replace" url:"replace,omitempty"`
	Search             bool     `json:"search" url:"search,omitempty"`
	Sentiment          bool     `json:"sentiment" url:"sentiment,omitempty"`
	SentimentThreshold float64  `json:"sentiment_threshold" url:"sentiment_threshold,omitempty"`
	SmartFormat        bool     `json:"smart_format" url:"smart_format,omitempty"`
	Start              string   `json:"start" url:"start,omitempty"`
	Summarize          bool     `json:"summarize" url:"summarize,omitempty"`
	Tag                []string `json:"tag" url:"tag,omitempty"`
	Translate          bool     `json:"translate" url:"translate,omitempty"`
	Utterances         bool     `json:"utterances" url:"utterances,omitempty"`
	UttSplit           bool     `json:"utt_split" url:"utt_split,omitempty"`
}

type UsageResponseDetail struct {
	Start    string  `json:"start"`
	End      string  `json:"end"`
	Hours    float64 `json:"hours"`
	Requests int     `json:"requests"`
}

type UsageSummary struct {
	Start      string                `json:"start"`
	End        string                `json:"end"`
	Resolution interface{}           `json:"resolution"`
	Results    []UsageResponseDetail `json:"results"`
}

type UsageRequestList struct {
	Page     int         `json:"page" url:"page,omitempty"`
	Limit    int         `json:"limit" url:"limit,omitempty"`
	Requests interface{} `json:"requests" url:"requests,omitempty"`
}

type UsageRequest struct {
	RequestId string      `json:"request_id" url:"request_id,omitempty"`
	Created   string      `json:"created" url:"created,omitempty"`
	Path      string      `json:"path" url:"path,omitempty"`
	Accessor  string      `json:"accessor" url:"accessor,omitempty"`
	Response  interface{} `json:"response" url:"response,omitempty"`
	Callback  interface{} `json:"callback" url:"callback,omitempty"`
}

func (dg *Client) ListRequests(projectId string, options UsageRequestListOptions) (UsageRequestList, error) {
	query, _ := query.Values(options)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/requests", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path, RawQuery: query.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
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

	var result UsageRequestList
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

func (dg *Client) GetRequest(projectId string, requestId string) (UsageRequest, error) {
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/requests/%s", dg.Path, projectId, requestId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path}
	req, err := http.NewRequest("GET", u.String(), nil)
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

	var result UsageRequest
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
		fmt.Printf("error getting request %s: %s\n", requestId, jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetFields(projectId string, options UsageRequestListOptions) (interface{}, error) {
	query, _ := query.Values(options)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/usage/fields", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path, RawQuery: query.Encode()}
	req, err := http.NewRequest("GET", u.String(), nil)
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

	var result interface{}
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
		fmt.Printf("error getting fields: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}

func (dg *Client) GetUsage(projectId string, options UsageOptions) (UsageSummary, error) {
	query, _ := query.Values(options)
	client := new(http.Client)
	path := fmt.Sprintf("%s/%s/usage", dg.Path, projectId)
	u := url.URL{Scheme: "https", Host: dg.Host, Path: path, RawQuery: query.Encode()}
	req, err := http.NewRequest("GET", u.String(), nil)
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

	var result UsageSummary
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
		fmt.Printf("error getting usage: %s\n", jsonErr.Error())
		return result, jsonErr
	} else {
		return result, nil
	}
}
