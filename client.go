package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Client struct {
	http    *http.Client
	baseURL string
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		http:    httpClient,
		baseURL: "https://api.spotify.com/v1/",
	}
}

func (c *Client) get(ctx context.Context, spotifyURL string, result interface{}) error {
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, spotifyURL, nil)
		if err != nil {
			return err
		}

		res, err := c.http.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusOK {
			return json.NewDecoder(res.Body).Decode(result)
		}
	}
}
