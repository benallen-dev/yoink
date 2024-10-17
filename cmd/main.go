package main

import (
	"os"

	"yoink/pkg/fourchan"
	"yoink/pkg/log"
)

func main() {
	logger := log.Default()
	foo, err := fourchan.GetPage("w", 1)
	if err != nil {
		logger.Error("Error fetching page", "error", err)
		os.Exit(1)
	}

	logger.Info("Got page 1", "threads", foo.ThreadCount())
}
