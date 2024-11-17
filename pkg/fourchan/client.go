package fourchan

import (
	"context"
	"net/http"
	"time"

	"yoink/pkg/cache"
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
	history    cache.TimeCache
}

func NewClient() Client {
	logger := log.Default()
	client := &http.Client{}
	history := *cache.NewCache("4chan")


	// Load history from disk
	if err := history.Load(); err != nil {
		logger.Error("Could not load cache from disk", "error", err)
	}
	
	// TODO: Add database of md5 hashes for images so we can avoid downloading
	//       images we already have.

	logger.Debug("Created new 4chan client")

	return Client{
		client:     client,
		limiter:    rate.NewLimiter(rate.Limit(1), 1),
		limiterCtx: context.Background(),
		history:    history,
	}
}

func (c *Client) Get(url string) (*http.Response, error) {
	logger := log.Default()
	logger.Info("GET (rate limited)", "url", url)

	req, err := http.NewRequest("GET", url, nil)

	if lastAccessed, ok := c.history.Get(url); ok {
		t := time.Unix(lastAccessed, 0).Format("Mon, 02 Jan 2006 15:04:05 GMT")
		logger.Debug("Setting If-Modified-Since header", "url", url, "lastAccessed", t)

		req.Header.Set("If-Modified-Since", t)
	}

	c.limiter.Wait(c.limiterCtx)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if (res.StatusCode == 200) { // Only update the last accessed time if the server response with 200
		c.history.Add(url)
	}
	return res, nil
}

func (c *Client) Close() {
	logger := log.Default()
	logger.Info("Closing 4chan client")

	if err := c.history.Persist(); err != nil {
		logger.Error("Could not persist cache", "error", err)
	}

	c.client.CloseIdleConnections()
}
