package main

import (
	"os"

	"github.com/charmbracelet/log"

	"yoink/pkg/fourchan"
)

// Just global because I'm heckin lazy
var Logger *log.Logger

func yoinkLogger() (out *log.Logger) {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.WarnLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.ErrorLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.DebugLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.FatalLevel].UnsetMaxWidth().Width(5)

	out = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		ReportTimestamp: true,
		Level: log.DebugLevel,
	})

	out.SetStyles(styles)

	return out
}

func main() {
	Logger = yoinkLogger()
	fourchan.Logger = Logger // wow that's nasty!

	foo, err := fourchan.GetPage("w", 1)
	if err != nil {
		Logger.Error("Error fetching page", "error", err)
		os.Exit(1)
	}

	Logger.Info("Got page 1", "threads", foo.ThreadCount())
}
