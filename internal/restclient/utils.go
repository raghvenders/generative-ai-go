package restclient

import "fmt"

const (
	Google_api_key    = "GOOGLE_API_KEY"
	Google_text_model = "PALM_API_MODEL"
	Palm2_api_url     = "GOOGLE_PAML2_URL"
	EMPTY             = ""

	//Experimental Beta values
	PALM2_GENERATIVE_URL = "https://generativelanguage.googleapis.com/v1beta2/models"
	GOOGLE_BISON_MODEL   = "text-bison-001"
)

func HyperLink(url string, text string) string {
	return fmt.Sprintf("  \x1b]8;;%s\x07%s\x1b]8;;\x07\u001b[0m", url, text)
}
