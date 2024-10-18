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
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.InfoLevel,
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

func Debug(msg string, keysAndValues ...interface{}) {
	logger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	logger.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	logger.Warn(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	logger.Error(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	logger.Fatal(msg, keysAndValues...)
}
