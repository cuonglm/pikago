package pikago

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	version       = "0.0.1"
	ua            = "pikago/" + version
	contentType   = "appilication/"
	defaultAPIURL = "https://pikab.in"
)

// PikabinClient is a HTTP client
type PikabinClient struct {
	client    *http.Client
	UserAgent string
	apiURL    string
}

// Option sets attributes for PikabinClient
type Option func(*PikabinClient) error

// Document represents content attributes
type Document struct {
	Content   string `json:"content"`
	Title     string `json:"title"`
	ExpiredAt string `json:"expired_at"`
	Syntax    string `json:"syntax"`
}

type payload struct {
	Payload Document `json:"document"`
}

// NewClient creates new PikabinClient
func NewClient(options ...Option) (*PikabinClient, error) {

	c := &PikabinClient{
		client:    http.DefaultClient,
		UserAgent: ua,
	}

	// Default API url
	err := APIUrl(defaultAPIURL)(c)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		err := option(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// APIUrl sets the API url option for gogi client
func APIUrl(url string) Option {
	return func(c *PikabinClient) error {
		c.apiURL = url
		return nil
	}
}

// Paste pastes a document to pikabin server
func (c *PikabinClient) Paste(d Document) (*http.Response, error) {
	payload := payload{Payload: d}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
