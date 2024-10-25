package fourchan

import (
	"context"
	"fmt"
	"path"
	"os"

	"yoink/pkg/log"
)

const QUEUE_LENGTH = 10_000

type QueueItem interface {
	getUrl() string
}

func ProcessQueue(ctx context.Context, q chan QueueItem) {
	logger := log.Default()

	for {
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
		if len(q) == 0 {
			return
		}
		logger.Debug(fmt.Sprintf("queue is now %d items long", len(q)))
	}
}

// I have somewhat overcomplicated this by using a chan instead of a slice
func NewQueue(board string) (q chan QueueItem) {
	logger := log.Default()
	
	// Make a directory to store results in
	dirPath := path.Join(getYoinkPath(), board)
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			logger.Error("Could not create dir", "path", dirPath, "error", err)
			return
		}
	}

	q = make(chan QueueItem, 10_000)
	q <- NewPageItem(board, 1)

	return q
}
