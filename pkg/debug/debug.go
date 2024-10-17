package debug

import (
	"encoding/json"
	"path"
	"os"

	"yoink/pkg/log"
)

func JsonToDisk(filename string, t any) {
	logger := log.Default()

	cwd, err := os.Getwd()
	if err != nil {
		logger.Error("Couldn't get cwd")
		return
	}

	debugPath := path.Join(cwd, "debug")
	_, err = os.Stat(debugPath)
	if os.IsNotExist(err) {
		os.Mkdir(debugPath, 0755)
	}

	// create file
	fullPath := path.Join(debugPath, filename+".json")
	f, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Could not create file", "path", fullPath)
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.Encode(t)
}

