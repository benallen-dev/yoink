package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"yoink/pkg/fourchan"
	"yoink/pkg/log"
)

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

	// Process all the things
	go func() {
		wg.Add(1)
		defer wg.Done()
		fourchan.ProcessQueue(fourCtx, fourchan.NewQueue("w"))
	}()

	// Block main thread until we receive an OS signal (exit)
	s := <-osSignal
	fmt.Println() // newline after ^C
	logger.Warn("Received OS signal", "signal", s)
	
	rootCancel()
	logger.Debug("Cancelled root context, waiting for wg.Wait()")
	wg.Wait()

	logger.Info("Exiting!")
}
