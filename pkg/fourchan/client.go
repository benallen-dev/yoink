package fourchan

import (
	"net/http"

	"yoink/pkg/log"

	"golang.org/x/time/rate"
)

// TODO: Implement a HTTP client that respects the API rules

// https://github.com/4chan/4chan-API/tree/master?tab=readme-ov-file
// 1. Do not make more than one request per second.
// 2. Thread updating should be set to a minimum of 10 seconds, preferably higher.
// 3. Use If-Modified-Since when doing your requests.
// 4. Make API requests using the same protocol as the app. Only use SSL when a user is accessing your app over HTTPS.

type Client struct {
	client *http.Client
	RateLimiter *rate.Limiter
}

func NewClient() Client {
	logger := log.Default()

	logger.Info("TODO: Finish implementing rate limited http client")

	return Client{
		client: http.DefaultClient,
		RateLimiter: rate.NewLimiter(rate.Limit(1), 1),
	}
}




