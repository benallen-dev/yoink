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
	"yoink/pkg/webui"
)

func main() {
	logger := log.Default()

	// What I really want to do is set up a system whereby in main all you do
	// is module.Start() and a thread that listens for shutdown and calls
	// context.Cancel(), there's way too much stuff going on in here

	// Set up listening for exit signals and gracefully exiting
	var wg sync.WaitGroup
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	// Set up contexts for each parallel process
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	// OS signal thread
	go func() {
		// Block thread until we receive an OS signal (exit)
		s := <-osSignal
		fmt.Println() // newline after ^C
		logger.Warn("Received OS signal", "signal", s)

		rootCancel()
		logger.Debug("Cancelled root context, waiting for wg.Wait()")
	}()

	// 4chan scraper threads
	wgCtx, wgCancel := context.WithCancel(rootCtx)
	defer wgCancel()

	wCtx, wCancel := context.WithCancel(rootCtx)
	defer wCancel()
	wg.Add(2) // If you do this in the goroutine, wg.Add is executed after wg.Wait

	// Process all the things
	go func() {
		defer wg.Done()
		fourchan.ProcessQueue(wgCtx, fourchan.NewQueue(wgCtx, "wg")) // wallpaper/general, not waitgroup
		logger.Info("Finished processing wg queue") 
	}()

	go func() {
		defer wg.Done()
		fourchan.ProcessQueue(wgCtx, fourchan.NewQueue(wCtx, "w")) // wallpaper/anime
		logger.Info("Finished processing w queue") 
	}()

	// Web ui thread
	httpCtx, httpCancel := context.WithCancel(rootCtx)
	defer httpCancel()
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info("Starting webui")
		webui.Listen(httpCtx)
		logger.Info("Halted webui")
	}()

	wg.Wait()

	logger.Info("Exiting!")
}
