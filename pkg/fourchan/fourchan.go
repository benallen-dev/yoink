package fourchan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func GetPage(board string, page int) (out Page, err error) {

	logger := *log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller: true,
		ReportTimestamp: true,
		TimeFormat: time.Kitchen,
	})
	url := fmt.Sprintf("https://a.4cdn.org/%s/%d.json", board, page)

	logger.Info("Fetching", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}
