package palm2

import (
	"errors"
	"os"

	"github.com/raghvenders/generative-ai-go/internal/restclient"
)

var (
	ErrMissingPalmAPIKey = errors.New("palm API key is missing. without that we cannot create Palm2 Client")
	ErrDummy             = errors.New("Dummy")
)

func NewPalm() (*restclient.Client, error) {

	token := os.Getenv(restclient.Google_api_key)

	if token == restclient.EMPTY {
		return nil, ErrMissingPalmAPIKey
	}

	return restclient.New("", "", "")
}
