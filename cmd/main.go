package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"yoink/pkg/log"
)

type QueueItem interface {
	getUrl() string
}

func main() {
	logger := log.Default()

	logger.Info("Starting the queue")

	// This is a bad idea so long as we only have 1 worker processing the
	// queue. A single queue item can push multiple new queue items and if
	// the queue fills up, it'll block until space becomes available - which
	// will never happen.
	q := make(chan QueueItem, 1000)
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	// seed the queue
	q <- PageQueueItem{
		board: "w",
		page:  1,
	}

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
			logger.Info("Received OS signal", "signal", s)
			return
		}
		logger.Info(fmt.Sprintf("q is now %d items long", len(q)))
	}

}
