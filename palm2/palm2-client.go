package palm2

import (
	"context"
	"errors"
	"log"

	"github.com/raghvenders/generative-ai-go/internal/restclient"
)

var (
	JoinErrs             = make([]error, 0, 1)
	ErrMissingPalmAPIKey = errors.New("ENV : PALM API key is missing. Kindly provide the key with env variable - GOOGLE_API_KEY")
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
	apiKey string
}

/*
func NewPalm() (Client, error) {

	token := os.Getenv(restclient.Google_api_key)

	if token == restclient.EMPTY {
		//return nil, ErrMissingPalmAPIKey
		JoinErrs = append(JoinErrs, ErrMissingPalmAPIKey)
	}

	url := os.Getenv(restclient.Palm2_api_url)

	if url == restclient.EMPTY {
		log.Printf("Default : Google PaLM2 Url is missing. Using default url : %s", restclient.HyperLink(restclient.PALM2_GENERATIVE_URL, restclient.PALM2_GENERATIVE_URL))
		url = restclient.PALM2_GENERATIVE_URL
	}

	model := os.Getenv(restclient.Google_text_model)

	if model == restclient.EMPTY {
		log.Printf("Default : Google PaLM2 Model is missing. Using default Model : %s", restclient.GOOGLE_BISON_MODEL)
		model = restclient.GOOGLE_BISON_MODEL
	}

	if len(JoinErrs) > 0 {
		return nil, func(e []error) error {
			var b []byte
			for i, err := range e {
				if i > 0 {
					b = append(b, '\n')
				}
				b = append(b, err.Error()...)
			}
			return fmt.Errorf(string(b))

		}(JoinErrs)
	}

	return restclient.New(url, token, model)
}
*/

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (rc *Client) ListModels() (*ModelsListResponse, error) {
	var err error

	listModels := &ModelsListResponse{}

	resErr := restclient.NewRestBuilder().
		WithBaseUrl("https://generativelanguage.googleapis.com/v1beta2/").
		WithPathParams("models").WithQueryParams(map[string]any{"key": rc.apiKey}).
		ResultJSON(listModels).
		ResultError(&err).Do(context.Background())

	if resErr != nil {
		return nil, resErr
	}

	log.Printf("%+v", *listModels)

	return listModels, nil
}
