package fourchan

import (
	"net/http"
	"time"

	"yoink/pkg/log"

	"golang.org/x/time/rate"
)

// https://github.com/4chan/4chan-API/tree/master?tab=readme-ov-file
// API Rules
// 1. Do not make more than one request per second.
// 2. Thread updating should be set to a minimum of 10 seconds, preferably higher.
// 3. Use If-Modified-Since when doing your requests.
// 4. Make API requests using the same protocol as the app. Only use SSL when a user is accessing your app over HTTPS.

type Client struct {
	client      *http.Client
	RateLimiter *rate.Limiter
}

func NewClient() Client {
	logger := log.Default()
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}

	logger.Info("TODO: Finish implementing rate limited http client")

	// TODO: Add rate limiting for 1 request per second

	// TODO: Add LRU cache for tracking the last time a route was accessed
	//		 so we can set If-Modified-Since headers.

	// TODO: Add database of md5 hashes for images so we can avoid downloading
	//       images we already have.

	return Client{
		client:      client,
		RateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}
