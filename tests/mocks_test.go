package deepgram_test

import (
	"bytes"
	"strings"
	"time"

	"github.com/deepgram-devs/deepgram-go-sdk/deepgram"
)

const (
	MockApiKey    = "testKey"
	MockApiSecret = "testSecret"
	MockApiDomain = "api.deepgram.test"
	MockEndPoint  = "https://api.deepgram.test/v1/listen"
	MockUuid      = "27e92bb2-8edc-4fdf-9a16-b56c78d39c5b"
	MockProjectId = MockUuid
	MockRequestId = MockUuid
	MockEmail     = "email@email.com"
	MockScope     = "read:mock"
	MockTag       = "string"
)

var MockDate = time.Now().Format(time.RFC3339)

var MockInvalidCredentials = map[string]any{
	"err_code":   "INVALID_AUTH",
	"err_msg":    "Invalid credentials.",
	"request_id": MockRequestId,
}

var MockKey = deepgram.Key{
	ApiKeyId: MockUuid,
	Key:      "string",
	Comment:  "string",
	Created:  "string",
	Scopes:   MockScopes,
}

var MockListKeys = deepgram.KeyResponse{
	ApiKeys: []deepgram.KeyResponseObj{
		{
			Member: MockMember,
			ApiKey: MockKey,
		},
	},
}

var MockMessageResponse = deepgram.Message{
	Message: "string",
}

var MockInvite = deepgram.InvitationOptions{
	Email: MockEmail,
	Scope: MockScope,
}

var MockInvites = deepgram.InvitationList{
	Invites: []deepgram.InvitationOptions{
		MockInvite,
		MockInvite,
	},
}

var MockScopes = []string{
	MockScope,
	MockScope,
}

var MockScopeList = deepgram.ScopeList{
	Scopes: MockScopes,
}

var MockTags = []string{MockTag}

var MockProjectKey = deepgram.KeyResponseObj{
	Member: deepgram.Member{
		MemberId:  MockUuid,
		Email:     MockEmail,
		FirstName: "string",
		LastName:  "string",
	},
	ApiKey: deepgram.Key{
		ApiKeyId: MockUuid,
		Key:      "string",
		Comment:  "string",
		Scopes:   MockScopes,
		Created:  MockDate,
	},
}

var MockBillingBalance = deepgram.Balance{
	BalanceId:       MockUuid,
	Amount:          0,
	Units:           "string",
	PurchaseOrderId: "string",
}

var MockBillingRequestList = deepgram.BalanceList{
	Balances: []deepgram.Balance{MockBillingBalance},
}

var MockUsgaeRequestListOptions = deepgram.UsageRequestListOptions{
	Start:  "string",
	End:    "string",
	Page:   1,
	Limit:  1,
	Status: "succeeded",
}

var MockUsageRequest = deepgram.UsageRequest{
	RequestId: "string",
	Created:   "string",
	Path:      "string",
	Accessor:  "string",
	Response:  MockMessageResponse,
}

var MockUsageRequestList = deepgram.UsageRequestList{
	Page:     1,
	Limit:    10,
	Requests: []deepgram.UsageRequest{MockUsageRequest},
}

var MockUsageOptions = deepgram.UsageOptions{
	Start:              "string",
	End:                "string",
	Accessor:           "string",
	Tag:                []string{MockTag},
	Method:             "sync",
	Model:              "string",
	Multichannel:       true,
	InterimResults:     true,
	Punctuate:          true,
	Ner:                true,
	Utterances:         true,
	Replace:            true,
	ProfanityFilter:    true,
	Keywords:           true,
	Sentiment:          true,
	Diarize:            true,
	DetectLanguage:     true,
	Search:             true,
	Redact:             true,
	Alternatives:       true,
	Numerals:           true,
	Numbers:            true,
	Translate:          true,
	DetectEntities:     true,
	DetectTopics:       true,
	Summarize:          true,
	Paragraphs:         true,
	UttSplit:           true,
	AnalyzeSentiment:   true,
	SmartFormat:        true,
	SentimentThreshold: 1.0,
}

var MockUsageResponseDetail = deepgram.UsageResponseDetail{
	Start:    "string",
	End:      "string",
	Hours:    1,
	Requests: 1,
}

var MockUsage = deepgram.UsageSummary{
	Start: "string",
	End:   "string",
	Resolution: map[string]any{
		"Units":  "string",
		"Amount": 1,
	},
	Results: []deepgram.UsageResponseDetail{
		MockUsageResponseDetail,
	},
}

var MockMember = deepgram.Member{
	MemberId:  MockUuid,
	Scopes:    MockScopes,
	Email:     MockEmail,
	FirstName: "string",
	LastName:  "string",
}

var MockMembers = deepgram.MemberList{
	Members: []deepgram.Member{MockMember},
}

var MockProject = deepgram.Project{
	ProjectId: MockUuid,
	Name:      "string",
	Company:   "string",
}

var MockPorjects = deepgram.ProjectResponse{
	Projects: []deepgram.Project{MockProject},
}

var MockProjectUpdate = deepgram.ProjectUpdateOptions{
	Name:    "string",
	Company: "string",
}

var MockStream = strings.NewReader("string")

var MockReadStreamSource = deepgram.ReadStreamSource{
	Stream:   MockStream,
	Mimetype: "video/mpeg",
}

var MockUrlSource = deepgram.UrlSource{
	Url: "string",
}

var MockBuffer = bytes.NewBufferString("string")

var MockBufferSource = deepgram.BufferSource{
	Buffer:   *MockBuffer,
	Mimetype: "video/mpeg",
}

var MockPrerecordedOptions = deepgram.PreRecordedTranscriptionOptions{
	Model:     "nova",
	Punctuate: true,
}

var MockMetaData = deepgram.Metadata{
	RequestId:      "string",
	TransactionKey: "string",
	Sha256:         "string",
	Created:        "string",
	Duration:       1.0,
	Channels:       1.0,
	ModelInfo: map[string]deepgram.ModelInfo{
		"b05e2505-2e49-4644-8e58-7878767ca60b": {
			Name:    "fake",
			Version: "version",
			Arch:    "arch",
		},
	},
	Models: []string{"string", "another string", "yet another string"},
}

var MockHit = deepgram.Hit{
	Confidence: 1.0,
	Start:      1.0,
	End:        1.0,
	Snippet:    "string",
}

var MockSearch = deepgram.Search{
	Query: "string",
	Hits:  []deepgram.Hit{MockHit},
}

var MockWordBase = deepgram.WordBase{
	Word:       "string",
	Start:      1.0,
	End:        1.0,
	Confidence: 1.0,
}

var MockAlternative = deepgram.Alternative{
	Transcript: "string",
	Confidence: 1.0,
	Words:      []deepgram.WordBase{MockWordBase},
}

var MockChannel = deepgram.Channel{
	Search:           []*deepgram.Search{&MockSearch},
	Alternatives:     []deepgram.Alternative{MockAlternative},
	DetectedLanguage: "string",
}

var MockPrerecordedResponse = deepgram.PreRecordedResponse{
	Request_id: MockUuid,
	Metadata:   MockMetaData,
	Results: deepgram.Results{
		Channels: []deepgram.Channel{MockChannel},
	},
}

const MockPrerecordedResponseJSON = `{
	"request_id": "27e92bb2-8edc-4fdf-9a16-b56c78d39c5b",
	"metadata": {
	  "transaction_key": "string",
	  "request_id": "string",
	  "sha256": "string",
	  "created": "string",
	  "duration": 1,
	  "channels": 1,
	  "models": [
		"string",
		"another string",
		"yet another string"
	  ],
	  "model_info": {
		"b05e2505-2e49-4644-8e58-7878767ca60b": {
		  "name": "fake",
		  "version": "version",
		  "arch": "arch"
		}
	  }
	},
	"results": {
	  "channels": [
		{
		  "search": [
			{
			  "query": "string",
			  "hits": [
				{
				  "confidence": 1,
				  "start": 1,
				  "end": 1,
				  "snippet": "string"
				}
			  ]
			}
		  ],
		  "alternatives": [
			{
			  "transcript": "string",
			  "confidence": 1,
			  "words": [
				{
				  "word": "string",
				  "start": 1,
				  "end": 1,
				  "confidence": 1
				}
			  ]
			}
		  ],
		  "detected_language": "string"
		}
	  ]
	}
  }`
