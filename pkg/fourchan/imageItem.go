package fourchan

import (
	"fmt"
	"io"
	"os"
	"path"

	"yoink/pkg/log"
)

// TODO: Move this to the correct place

type ImageItem struct {
	board    string
	filename string
}

func (i ImageItem) getUrl() string {
	return fmt.Sprintf("https://i.4cdn.org/%s/%s", i.board, i.filename)
}

func handleImageQueueItem(i ImageItem) {
	logger := log.Default()

	url := i.getUrl()
	logger.Debug("Fetching", "url", url)

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

	basePath, err := os.Getwd()
	if err != nil {
		logger.Error("Couldn't get cwd")
		return
	}

	// check if dir exists, create it if not
	dirPath := path.Join(basePath, "img", i.board)
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		os.Mkdir(dirPath, 0755)
	}

	// create file
	fullPath := path.Join(basePath, i.board, i.filename)
	f, err := os.Create(fullPath)
	if err != nil {
		logger.Error("Could not create file", "path", fullPath)
		return
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		logger.Error("Could not write file", "path", fullPath)
		return
	}
}
