package fourchan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"yoink/pkg/log"
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

	var thread Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		logger.Error("Could not decode thread", "error", err, "url", url, "board", i.board, "op", i.op)
	}

	for _, p := range thread.Posts {
		if p.Filename != "" && p.W == 3840 && p.H == 2160{

			q <- ImageItem{
				board:    i.board,
				filename: strconv.Itoa(p.Tim) + p.Ext,
			}
		}
	}
}
