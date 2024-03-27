package astrometry

import (
	"net/http"
	"time"
)

// Constants

const actualAPI = "https://nova.astrometry.com/api"

// Client

type Client struct {
	apiKey      string
	BaseURL     string
	SessionKey  string
	Submissions []string
	httpClient  *http.Client
}

func NewAstrometryClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		BaseURL: actualAPI,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Methods

func (c *Client) Connect() {

}

func (c *Client) UploadFile(file string) string {
	return ""
}

func (c *Client) ReviewFile(subID string) {

}
