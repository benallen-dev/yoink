package main

import (
	"context"
	"fmt"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"

	"yoink/pkg/config"
	"yoink/pkg/fourchan"
	"yoink/pkg/log"
	"yoink/pkg/webui"
)

func worker(id int, filechan <-chan os.DirEntry) {

	dl := log.Default()
	logger := dl.With("worker", id)

	// get width and height
	// if not 3840x2160, delete
	for file := range filechan {

		e := path.Ext(file.Name())
		if e == ".jpg" || e == ".png" || e == ".jpeg" {
			//logger.Info("Found image", "file", file.Name())
			
			file, err := os.Open(path.Join(config.NewDir, file.Name()))
			if err != nil {
				logger.Error("Could not open file", "error", err)
				continue
			}

			img, _, err := image.Decode(file)
			if err != nil {
				logger.Error("Could not decode image", "error", err)
				err := os.Rename(path.Join(config.NewDir, file.Name()), path.Join(config.FaultyDir, path.Base(file.Name())))
				if err != nil {
					logger.Error("Could not move file to faulty", "error", err)
					continue
				}
			}

			if img.Bounds().Dx() != 3840 || img.Bounds().Dy() != 2160 {
				logger.Info("Deleting image", "file", file.Name())
				err := os.Rename(path.Join(file.Name()), path.Join(config.DeletedDir, path.Base(file.Name())))
				if err != nil {
					logger.Error("Could not move file to deleted", "error", err)
					continue
				}

			}
		}
	}
}

func cleanup() {
	logger := log.Default()

	logger.Info("Cleaning up")
	logger.Info("New dir", "path", config.NewDir)
	// Let's tidy up all those weird sizes
	files, err := os.ReadDir(config.NewDir)
	if err != nil {
		logger.Error("Could not read new dir", "error", err)
		return
	}
	// create 4 workers
	jobs := make(chan os.DirEntry)
	for i := 0; i <= 4; i++ {
		go worker(i, jobs)
	}

	for _, file := range files {
		// get width and height
		jobs <- file
	}
	close(jobs)

}

func main() {
	logger := log.Default()
	cleanup()

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
