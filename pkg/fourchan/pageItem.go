package fourchan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"yoink/pkg/log"
)

type PageQueueItem struct {
	board string
	page  int
}

func NewPageQueueItem(board string, page int) PageQueueItem {
	return PageQueueItem{
		board: board,
		page:  page,
	}
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

	var page Page
	json.NewDecoder(resp.Body).Decode(&page)

	for _, t := range page.Threads {
		first := t.Posts[0]

		q <- ThreadQueueItem{
			board: i.board,
			op:    first.No,
		}
	}

}
