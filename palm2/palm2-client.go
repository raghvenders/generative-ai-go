package palm2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/raghvenders/generative-ai-go/internal/restclient"
)

type HarmCategory int
type HarmBlockThreshold int
type HarmProbability int
type BlockedReason int

const (
	HARM_CATEGORY_UNSPECIFIED HarmCategory = iota
	HARM_CATEGORY_DEROGATORY
	HARM_CATEGORY_TOXICITY
	HARM_CATEGORY_VIOLENCE
	HARM_CATEGORY_SEXUAL
	HARM_CATEGORY_MEDICAL
	HARM_CATEGORY_DANGEROUS

	HARM_BLOCK_THRESHOLD_UNSPECIFIED HarmBlockThreshold = iota
	BLOCK_LOW_AND_ABOVE
	BLOCK_MEDIUM_AND_ABOVE
	BLOCK_ONLY_HIGH
	BLOCK_NONE

	HARM_PROBABILITY_UNSPECIFIED HarmProbability = iota
	NEGLIGIBLE
	LOW
	MEDIUM
	HIGH

	BLOCKED_REASON_UNSPECIFIED BlockedReason = iota
	SAFETY
	OTHER
)

var (
	JoinErrs             = &MultipleErrors{errs: make([]error, 0, 1)}
	ErrMissingPalmAPIKey = errors.New("PaLM2 API Key is missing as env variable - GOOGLE_API_KEY. Kindly Provide a Valid key to Access PaLM2 API")
)

type GenerativeModels struct {
	Name                       string   `json:"name"`
	BaseModelId                string   `json:"baseModelId"`
	Version                    string   `json:"version"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description"`
	InputTokenLimit            int64    `json:"inputTokenLimit"`
	OutputTokenLimit           int64    `json:"outputTokenLimit"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	Temperature                float64  `json:"temperature"`
	TopP                       float64  `json:"topP"`
	TopK                       int64    `json:"topK"`
}

type ModelsListResponse struct {
	Models        []GenerativeModels `json:"models"`
	NextPageToken string             `json:"nextPageToken"`
}

type SafetySetting struct {
	Category  HarmCategory
	Threshold HarmBlockThreshold
}

type SafetyRating struct {
	Category    HarmCategory
	Probability HarmProbability
}

type TextPrompt struct {
	Text string `json:"text"`
}

type PalmTextRequest struct {
	Prompt          TextPrompt      `json:"prompt"`
	SafetySettings  []SafetySetting `json:"safetySettings"`
	StopSequences   []string        `json:"stopSequences"`
	Temperature     float32         `json:"temperature"`
	CandidateCount  uint8           `json:"candidateCount"`
	MaxOutputTokens uint32          `json:"maxOutputTokens"`
	TopP            float32         `json:"topP"`
	TopK            uint32          `json:"topK"`
}

type PalmTextResponse struct {
	Candidates     []TextCompletion `json:"candidates"`
	Filters        []ContentFilter  `json:"filters"`
	SafetyFeedback SafetyFeedback   `json:"safetyFeedback"`
}

type SafetyFeedback struct {
	Rating  SafetyRating  `json:"rating"`
	Setting SafetySetting `json:"setting"`
}

type ContentFilter struct {
	Reason  BlockedReason `json:"Reason"`
	Message string        `json:"message"`
}

type TextCompletion struct {
	Output           string          `json:"Output"`
	SafetyRatings    []SafetyRating  `json:"safetyRatings"`
	CitationMetadata CitationMetdata `json:"citationMetadata"`
}

type CitationMetdata struct {
	CitationSources []CitationSource `json:"citationSources"`
}

type CitationSource struct {
	StartIndex uint32 `json:"startIndex"`
	EndIndex   uint32 `json:"endIndex"`
	Uri        string `json:"uri"`
	License    string `json:"license"`
}
type Client struct {
	apiKey  string
	baseURL string
}

type MultipleErrors struct {
	errs []error
}

func (e *MultipleErrors) Error() string {
	var s strings.Builder
	for _, err := range e.errs {
		s.WriteString(err.Error())
	}
	return s.String()
}

func (e *MultipleErrors) Unwrap() []error {
	return e.errs
}

func NewPalm() (*Client, error) {

	var token string
	if os.Getenv(restclient.Google_api_key) != restclient.EMPTY {
		token = os.Getenv(restclient.Google_api_key)
	}

	if token == restclient.EMPTY {
		JoinErrs.errs = append(JoinErrs.errs, ErrMissingPalmAPIKey)
	}

	url := os.Getenv(restclient.Palm2_api_url)

	if url == restclient.EMPTY {
		log.Printf("Default : Google PaLM2 Url is missing. Using default url : %v", restclient.HyperLink(restclient.PALM2_GENERATIVE_URL, restclient.PALM2_GENERATIVE_URL))
		url = restclient.PALM2_GENERATIVE_URL
	}

	if len(JoinErrs.errs) > 0 {
		return nil, JoinErrs
	}

	return &Client{apiKey: token, baseURL: url}, nil
}

func (rc *Client) ListModels(queryParams ...string) (*ModelsListResponse, error) {
	var err error

	listModels := &ModelsListResponse{}

	var params map[string]any

	if len(queryParams)%2 != 0 {
		log.Printf("Ignoring Query Params as it has to be >= 2 and should be even :  %d", len(queryParams))
	} else {
		params := map[string]any{}
		for i := 0; i < len(queryParams); i += 2 {
			params[queryParams[i]] = queryParams[i+1]
		}
	}

	resErr := restclient.NewRestBuilder().
		WithBaseUrl(rc.baseURL).
		WithHeaders(map[string][]string{"x-goog-api-key": {rc.apiKey}}).
		WithPathParams("models").
		WithQueryParams(params).
		ResultJSON(listModels).
		ResultError(&err).Do(context.Background(), http.MethodGet)

	if resErr != nil {
		return nil, resErr
	}

	return listModels, nil
}

func (rc *Client) GenerateText(prompt string, queryParams ...string) (any, error) {
	var err error

	var rawResponse string

	var params map[string]any

	if len(queryParams)%2 != 0 {
		log.Printf("Ignoring Query Params as it has to be >= 2 and should be even :  %d", len(queryParams))
	} else {
		params := map[string]any{}
		for i := 0; i < len(queryParams); i += 2 {
			params[queryParams[i]] = queryParams[i+1]
		}
	}

	resErr := restclient.NewRestBuilder().
		WithBaseUrl(rc.baseURL).
		WithHeaders(map[string][]string{"x-goog-api-key": {rc.apiKey}, "Content-Type": {"application/json"}}).
		WithRawBody(fmt.Sprintf(`{
			"prompt": {
			"text": "%s"
		}
	}`, prompt)).
		WithPathParams("models/%s/:%s", "text-bison-001", "generateText").
		WithQueryParams(params).
		ResultRaw(&rawResponse).
		ResultError(&err).Do(context.Background(), http.MethodPost)

	if resErr != nil {
		return nil, resErr
	}

	return rawResponse, nil
}
