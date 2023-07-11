package palm2

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/raghvenders/generative-ai-go/internal/restclient"
)

var (
	JoinErrs             = &MultipleErrors{errs: make([]error, 0, 1)}
	ErrMissingPalmAPIKey = errors.New("A")
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

type Client struct {
	apiKey  string
	model   string
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

/*
func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}
*/

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
		ResultError(&err).Do(context.Background())

	if resErr != nil {
		return nil, resErr
	}

	return listModels, nil
}
