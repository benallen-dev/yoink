package fourchan

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"yoink/pkg/config"
	"yoink/pkg/log"
)

func imageOk(p Post) bool {
	aspect := float64(p.W) / float64(p.H)
	aspectOk := aspect < 1.8 && aspect > 1.77 // 16:9 is 1.777...
	widthOk := p.W >= 3840
	heightOk := p.H >= 2160

	return p.Filename != "" && widthOk && heightOk && aspectOk
}

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

	switch resp.StatusCode {
	case 200:
		// continue
	case 304:
		logger.Debugf("Got non-OK statuscode %s for %s", resp.Status, url)
		return
	default:
		logger.Warnf("Got non-OK statuscode %s for %s", resp.Status, url)
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

		//if p.Filename != "" && p.W >= 3840 && p.H >= 2160 { 
		if imageOk(p) {
			tim := strconv.Itoa(p.Tim)

			q <- ImageItem{
				board:    i.board,
				filename: tim + p.Ext,
				tim:      tim,
			}
		}
	}
}
