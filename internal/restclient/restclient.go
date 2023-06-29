package restclient

// Client is a client for the OpenAI API.
type Client struct {
	BaseURL string
	Model   string
	AuthKey string
}

// New returns a new OpenAI client.
func New(baseURL string, token string, model string) (*Client, error) {
	c := &Client{
		Model:   model,
		BaseURL: baseURL,
		AuthKey: token,
	}

	return c, nil
}
