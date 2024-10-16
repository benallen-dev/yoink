package main

import (
	"os"
	"os/signal"
	"syscall"

	"yoink/pkg/log"
	"yoink/pkg/fourchan"
)

func main() {
	logger := log.Default()

	logger.Info("Starting the queue")

	// This is a bad idea so long as we only have 1 worker processing the
	// queue. A single queue item can push multiple new queue items and if
	// the queue fills up, it'll block until space becomes available - which
	// will never happen.
	q := make(chan fourchan.QueueItem, 1000)
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	// seed the queue
	q <- fourchan.NewPageQueueItem("w", 1)

	// kinda meh to implicitly block but whatever let's just see if this works
	fourchan.ProcessQueue(q, osSignal)
}
