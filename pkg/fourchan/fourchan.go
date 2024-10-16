package fourchan

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func GetPage(board string, page int) (out Page, err error) {

	url := fmt.Sprintf("https://a.4cdn.org/%s/%d.json", board, page)

	Logger.Info("Fetching", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&out)
	return out, nil
}
