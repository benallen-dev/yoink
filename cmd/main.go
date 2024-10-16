package main

import (
	"os"

	"github.com/charmbracelet/log"

	"yoink/pkg/fourchan"
)

// Just global because I'm heckin lazy
var logger log.Logger

func customLogger() (logger log.Logger) {
	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel].SetString("ERROR")
	styles.Levels[log.DebugLevel].SetString("DEBUG")
	styles.Levels[log.FatalLevel].SetString("FATAL")

	logger = *log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		ReportTimestamp: true,
	})

	logger.SetStyles(styles)

	return logger
}

func main() {
	logger = customLogger()

	foo, err := fourchan.GetPage("w", 1)
	if err != nil {
		logger.Error("Error fetching page", "error", err)
		os.Exit(1)
	}

	logger.Info("Got page 1", "threads", foo.ThreadCount())

	logger.Debug("I'm a debug message")
	logger.Warn("I'm a warning")
	logger.Error("I'm an error")
	logger.Fatal("I'm fatal")

}
