package deepgram

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
