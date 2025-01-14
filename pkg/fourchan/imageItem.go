package fourchan

import (
	"fmt"
	"io"
	"os"
	"path"

	"yoink/pkg/config"
	"yoink/pkg/log"
)

type ImageItem struct {
	board    string
	filename string
	tim      string
}

func (i ImageItem) getUrl() string {
	return fmt.Sprintf("https://i.4cdn.org/%s/%s", i.board, i.filename)
}

func handleImageQueueItem(i ImageItem) {
	logger := log.Default()

	url := i.getUrl()
	logger.Debug("Fetching", "url", url)

	fullPath := path.Join(config.DataPath(), "new", i.filename)
	stat, err := os.Stat(fullPath)
	if err == nil && stat.Size() > 0 {
		logger.Debug("File already exists", "path", fullPath)
		return
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Warn("Could not fetch image", "url", url, "board", i.board)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Warn(fmt.Sprintf("Got non-OK statuscode %s for %s", resp.Status, url))
		return
	}

	// create file
	f, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Could not create file", "path", fullPath, "error", err)
		return
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		logger.Error("Could not write file", "path", fullPath, "error", err)
		return
	}

	logger.Info("Downloaded", "path", fullPath)
}
