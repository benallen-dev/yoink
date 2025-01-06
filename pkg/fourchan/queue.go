package fourchan

import (
	"context"
	"fmt"

	"yoink/pkg/log"
)

const QUEUE_LENGTH = 10_000

type QueueItem interface {
	getUrl() string
}

func ProcessQueue(ctx context.Context, q chan QueueItem) {
	logger := log.Default()
	logger = *logger.With("module", "fourchan")

	for {
		select {
		case <-ctx.Done():
			httpClient.Close()
			logger.Info("Received context.Done, stopping ProcessQueue")
			return
		case i := <-q:
			switch i.(type) {
			case PageItem:
				handlePageQueueItem(i.(PageItem), q)
			case ThreadItem:
				handleThreadQueueItem(i.(ThreadItem), q)
			case ImageItem:
				handleImageQueueItem(i.(ImageItem))
			}
		}

		if len(q) == 0 {
			httpClient.Close()
			logger.Info("Queue is empty, exiting ProcessQueue")
			return
		}

		if len(q)%50 == 0 {
			logger.Info(fmt.Sprintf("queue is now %d items long", len(q)))
		} else {
			logger.Debug(fmt.Sprintf("queue is now %d items long", len(q)))
		}
	}
}

// I have somewhat overcomplicated this by using a chan instead of a slice
func NewQueue(ctx context.Context, board string) (q chan QueueItem) {
	q = make(chan QueueItem, 10_000)
	q <- NewPageItem(board, 1)
	return q
}
