package sleuth

import (
	"net/http"
)

type Client struct {
	httpClient       *http.Client
	organizationSlug string
	apiKey           string
}

type NewClientOption func(*Client)

func NewClient(
	organizationSlug string,
	apiKey string,
	options ...NewClientOption,
) *Client {
	client := &Client{
		organizationSlug: organizationSlug,
		apiKey:           apiKey,
	}

	for _, option := range options {
		option(client)
	}

	return client
}
