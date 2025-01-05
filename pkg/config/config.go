package config

import (
	"os"
	"path"

	"yoink/pkg/log"
)

var (
	NewDir     = path.Join(DataPath(), "new")
	FaultyDir  = path.Join(DataPath(), "faulty")
	KeepDir    = path.Join(DataPath(), "categorised", "keep")
	DiscardDir = path.Join(DataPath(), "categorised", "discard")
	AnimeDir   = path.Join(DataPath(), "categorised", "anime-nsfw")
	NsfwDir    = path.Join(DataPath(), "categorised", "nsfw")
	DeletedDir = path.Join(DataPath(), "deleted")

	ImageW = 3840
	ImageH = 2160
)

func init() {
	logger := log.Default()
	// Ensure directories exist
	dirs := []string{NewDir, KeepDir, DiscardDir, AnimeDir, NsfwDir, FaultyDir, DeletedDir}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Fatal("Could not create directory", "dir", dir, "error", err)
		}
	}
}

func DataPath() (out string) {
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
