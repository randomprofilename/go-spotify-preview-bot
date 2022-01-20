package spotify_api

import "time"

type Client struct {
	clientId, clientSecret string

	token struct {
		expires *time.Time
		value   string
	}
}

func NewClient(clientId, clientSecret string) *Client {
	return &Client{
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}
