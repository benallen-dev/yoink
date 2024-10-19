package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"yoink/pkg/fourchan"
	"yoink/pkg/log"
)

func secondGoRoutine(ctx context.Context) {
	logger := log.Default()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Second goroutine received context.Done, returning")
			return
		default:
			logger.Info("Second goroutine is also still running")
			time.Sleep(1 * time.Second)
		}
	}
}

func main() {
	logger := log.Default()

	// Set up listening for exit signals and gracefully exiting
	var wg sync.WaitGroup
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	// Set up contexts for each parallel process
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	fourCtx, fourCancel := context.WithCancel(rootCtx)
	defer fourCancel()

	fooCtx, fooCancel := context.WithCancel(rootCtx)
	defer fooCancel()

	// Process all the things
	go func() {
		wg.Add(1)
		defer wg.Done()
		fourchan.ProcessQueue(fourCtx, fourchan.NewQueue("w"))
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		secondGoRoutine(fooCtx)
	}()


	// Block main thread until we receive an OS signal (exit)
	s := <-osSignal
	logger.Info("Received OS signal", "signal", s)
	rootCancel()
	logger.Info("Cancelled root context, waiting for wg.Wait()")
	wg.Wait()
	logger.Info("Exiting")
}
