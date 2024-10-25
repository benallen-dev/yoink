package fourchan

import (
	"os"
	"path"

	"yoink/pkg/log"
)

func getYoinkPath() (out string) {
	logger := log.Default()
	out = os.Getenv("YOINK_BASE_PATH")

	if out == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			logger.Error("Could not get homedir", "error", err)
			return ""
		}
		out = path.Join(homedir, "yoink")
	}

	logger.Debug("Base path", "path", out)

	// check if dir exists, create it if not
	_, err := os.Stat(out)
	if os.IsNotExist(err) {
		err = os.Mkdir(out, 0755)
		if err != nil {
			logger.Error("Could not create dir", "path", out, "error", err)
			return
		}
	}
	
	return out
}
