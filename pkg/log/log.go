package log

import (
	"os"

	"github.com/charmbracelet/log"
)

var logger *log.Logger

func init() {
	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.WarnLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.ErrorLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.DebugLevel].UnsetMaxWidth().Width(5)
	styles.Levels[log.FatalLevel].UnsetMaxWidth().Width(5)

	logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		ReportTimestamp: true,
		Level: log.DebugLevel,
	})

	logger.SetStyles(styles)
}

// Global returns the global logger, configured with sensible defaults.
func Default() log.Logger {
	return *logger
}

func Custom(options log.Options) *log.Logger {
	return log.NewWithOptions(os.Stderr, options)
}
