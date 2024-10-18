package fourchan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	//"yoink/pkg/debug"
	"yoink/pkg/log"
)

type ThreadQueueItem struct {
	board string
	op    int
}

func (i ThreadQueueItem) getUrl() string {
	return fmt.Sprintf("https://a.4cdn.org/%s/thread/%d.json", i.board, i.op)
}

func handleThreadQueueItem(i ThreadQueueItem, q chan QueueItem) {
	logger := log.Default()

	url := i.getUrl()
	logger.Info("Fetching", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Warn("Could not fetch thread", "url", url, "board", i.board, "op", i.op)
		return
	}
	defer resp.Body.Close()

	var thread Thread
	err = json.NewDecoder(resp.Body).Decode(&thread)
	if err != nil {
		logger.Error("Could not decode thread", "error", err, "url", url, "board", i.board, "op", i.op)
	}

	//debug.JsonToDisk(fmt.Sprintf("%d", i.op), thread)

	for _, p := range thread.Posts {
		if p.Filename != "" {

			q <- ImageQueueItem{
				board:    i.board,
				filename: strconv.Itoa(p.Tim) + p.Ext,
			}
		}
	}
}
