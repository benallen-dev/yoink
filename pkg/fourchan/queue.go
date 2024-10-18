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
			case PageQueueItem:
				handlePageQueueItem(i.(PageQueueItem), q)
			case ThreadQueueItem:
				handleThreadQueueItem(i.(ThreadQueueItem), q)
			case ImageQueueItem:
				handleImageQueueItem(i.(ImageQueueItem))
			}
		case s := <-osSignal:
			logger.Info("Received OS signal, stopping ProcessQueue", "signal", s)
			return
		}
		logger.Debug(fmt.Sprintf("queue is now %d items long", len(q)))
	}
}

