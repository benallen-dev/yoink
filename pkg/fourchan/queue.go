package fourchan

import (
	"context"
	"fmt"

	"yoink/pkg/log"
)

type QueueItem interface {
	getUrl() string
}

func ProcessQueue(ctx context.Context, q chan QueueItem) {
	logger := log.Default()

	for len(q) > 0 {
		select {
		case <-ctx.Done():
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
		logger.Debug(fmt.Sprintf("queue is now %d items long", len(q)))
	}
}

func NewQueue(board string) (q chan QueueItem) {
	q = make(chan QueueItem, 1000)
	q <- NewPageItem(board, 1)

	return q
}
