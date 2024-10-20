package fourchan

import (
	"context"
	"net/http"

	"yoink/pkg/log"

	"golang.org/x/time/rate"
)

// https://github.com/4chan/4chan-API/tree/master?tab=readme-ov-file
// API Rules
// 1. Do not make more than one request per second.
// 2. Thread updating should be set to a minimum of 10 seconds, preferably higher.
// 3. Use If-Modified-Since when doing your requests.
// 4. Make API requests using the same protocol as the app. Only use SSL when a user is accessing your app over HTTPS.

// Package-level httpClient for making requests to 4chan API
var httpClient = NewClient()

type Client struct {
	client     *http.Client
	limiter    *rate.Limiter
	limiterCtx context.Context
}

func NewClient() Client {
	logger := log.Default()
	client := &http.Client{}

	// TODO: Add persistent LRU cache for tracking the last time a route was accessed
	//		 so we can set If-Modified-Since headers.

	// TODO: Add database of md5 hashes for images so we can avoid downloading
	//       images we already have.

	logger.Debug("Created new 4chan client")

	return Client{
		client:     client,
		limiter:    rate.NewLimiter(rate.Limit(1), 1),
		limiterCtx: context.Background(),
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	logger := log.Default()
	logger.Info("GET (rate limited)", "url", url)
	c.limiter.Wait(c.limiterCtx)
	return c.client.Get(url)
}
