package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yoink/pkg/fourchan"
	"yoink/pkg/log"
)

// TODO: Move this to the correct place

type PageQueueItem struct {
	board string
	page  int
}

func (i PageQueueItem) getUrl() string {
	return fmt.Sprintf("https://a.4cdn.org/%s/%d.json", i.board, i.page)
}

func handlePageQueueItem(i PageQueueItem, q chan QueueItem) {
	logger := log.Default()
	url := i.getUrl()

	logger.Info("Fetching", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		// just exit, I suppose
		log.Warn("Could not fetch page", "url", url, "board", i.board, "page", i.page)
		return
	}
	defer resp.Body.Close()

	var page fourchan.Page
	json.NewDecoder(resp.Body).Decode(&page)

	for _, t := range page.Threads {
		first := t.Posts[0]
		q <- ThreadQueueItem{
			board: i.board,
			op: first.No,
		}
	}

}
