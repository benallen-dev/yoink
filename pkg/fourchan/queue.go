package fourchan

import (
	"fmt"
	"os"

	"yoink/pkg/log"
)

type QueueItem interface {
	getUrl() string
}

func ProcessQueue(q chan QueueItem, osSignal chan os.Signal) {
	logger := log.Default()

	for len(q) > 0 {
		select {
		case i := <-q:
			switch i.(type) {
			case PageItem:
				handlePageQueueItem(i.(PageItem), q)
			case ThreadItem:
				handleThreadQueueItem(i.(ThreadItem), q)
			case ImageItem:
				handleImageQueueItem(i.(ImageItem))
			}
		case s := <-osSignal:
			logger.Info("Received OS signal, stopping ProcessQueue", "signal", s)
			return
		}
		logger.Debug(fmt.Sprintf("queue is now %d items long", len(q)))
	}
}

func NewQueue(board string) (q chan QueueItem) {
	q = make(chan QueueItem, 1000)
	q <- NewPageItem(board,1)

	return q
}
