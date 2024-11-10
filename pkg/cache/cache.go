package cache

import (
	// "context"
	"fmt"
	// "path"
	// "os"

	"os"
	"yoink/pkg/log"
)

type Cache struct {
	name   string
	images map[string]bool
}

func NewCache() Cache {
	logger := log.Default()

	logger.Debug("Created new cache")
	return Cache{
		images: map[string]bool{
			"1234": true,
			"5678": true,
		},
	}
}

func (c Cache) Contains(key string) bool {
	_, ok := c.images[key]
	return ok
}

func (c Cache) Add(key string) {
	c.images[key] = true
}

func (c Cache) Persist() error {
	logger := log.Default()
	logger.Debug("Persisting cache")

	// I need a smart encoding scheme for this, so I'm just gonna yolo a list of /r/n separated keys
	filecontents := ""
	for key := range c.images {
		filecontents += key + "\n"
	}

	// A smarter way to do this once the file gets large is to append after add
	f, err := os.Create("cachefile")
	if err != nil {
		return err
	}

	_, err = f.WriteString(filecontents)
	if err != nil {
		return err
	}

	return nil
	 
}

func (c Cache) String() string {
	out := fmt.Sprintf("Cache: %s", c.name)

	first := true
	for key := range c.images {
		if !first {
			out += ", "
		} else {
			first = false
		}
		out += key
	}

	return out
}
