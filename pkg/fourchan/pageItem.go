package fourchan

import (
	"encoding/json"
	"fmt"

	"yoink/pkg/log"
)

type PageItem struct {
	board string
	page  int
}

func NewPageItem(board string, page int) PageItem {
	return PageItem{
		board: board,
		page:  page,
	}
}

func (i PageItem) getUrl() string {
	return fmt.Sprintf("https://a.4cdn.org/%s/%d.json", i.board, i.page)
}

func handlePageQueueItem(i PageItem, q chan QueueItem) {
	logger := log.Default()
	url := i.getUrl()

	logger.Debug("Fetching", "url", url)
	resp, err := httpClient.Get(url)
	if err != nil {
		logger.Warn("Could not fetch page", "url", url, "board", i.board, "page", i.page)
		return
	}
	defer resp.Body.Close()

	var page Page
	json.NewDecoder(resp.Body).Decode(&page)

	for _, t := range page.Threads {
		first := t.Posts[0]

		q <- ThreadItem{
			board: i.board,
			op:    first.No,
		}
	}

}
