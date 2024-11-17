package fourchan

import (
	"context"
	"fmt"
	"os"
	"path"

	"yoink/pkg/config"
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
	logger := log.Default()

	// Make a directory to store results in
	dirPath := path.Join(config.DataPath(), board)
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
	q <- NewPageItem(board, 2)
	q <- NewPageItem(board, 3)
	q <- NewPageItem(board, 4)
	q <- NewPageItem(board, 5)
	q <- NewPageItem(board, 6)
	q <- NewPageItem(board, 7)
	q <- NewPageItem(board, 8)
	q <- NewPageItem(board, 9)
	q <- NewPageItem(board, 10)
	q <- NewPageItem(board, 11)
	q <- NewPageItem(board, 12)
	q <- NewPageItem(board, 13)
	q <- NewPageItem(board, 14)
	q <- NewPageItem(board, 15)
	return q
}
