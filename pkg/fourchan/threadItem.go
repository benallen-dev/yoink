package fourchan

import (
	"os"
	"path"
	"encoding/json"
	"fmt"
	"strconv"

	"yoink/pkg/log"
	"yoink/pkg/config"
)

type ThreadItem struct {
	board string
	op    int
}

func (i ThreadItem) getUrl() string {
	return fmt.Sprintf("https://a.4cdn.org/%s/thread/%d.json", i.board, i.op)
}

func handleThreadQueueItem(i ThreadItem, q chan QueueItem) {
	logger := log.Default()
	
	url := i.getUrl()
	logger.Debug("Fetching", "url", url)

	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Warn("Could not fetch thread", "url", url, "board", i.board, "op", i.op)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.Warn(fmt.Sprintf("Got non-OK statuscode %s for %s", resp.Status, url))
		return
	}

	var thread Thread
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		logger.Error("Could not decode thread", "error", err, "url", url, "board", i.board, "op", i.op)
	}

	for _, p := range thread.Posts {

		fullPath := path.Join(config.DataPath(), i.board, p.Filename)
		stat, err := os.Stat(fullPath)
		if err == nil && stat.Size() > 0 {
			logger.Debug("File already exists", "path", fullPath)
			continue
		}

		if p.Filename != "" && p.W == 3840 && p.H == 2160{
			tim := strconv.Itoa(p.Tim)

			q <- ImageItem {
				board:    i.board,
				filename: tim + p.Ext,
				tim: tim,
			}
		}
	}
}
