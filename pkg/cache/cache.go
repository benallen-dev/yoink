package cache

import (
	"fmt"
	"time"

	"bytes"
	"encoding/gob"
	"os"
	"path"

	"yoink/pkg/config"
	"yoink/pkg/log"
)


func init() {
	gob.Register(&TimeCache{})
}

// TimeCache is a cache that stores keys and their last access time
// Times are stored internally as Unix timestamps as an int64
type TimeCache struct {
	name    string
	entries map[string]int64
}

func (c TimeCache) String() string {
	out := fmt.Sprintf("Cache: %s\n", c.name)

	for k, v := range c.entries {
		out += fmt.Sprintf("%s: %d\n", k, v)
	}

	return out
}

func NewCache(name string) *TimeCache {
	return &TimeCache{
		name:    name,
		entries: make(map[string]int64),
	}
}

func (c TimeCache) Contains(key string) bool {
	_, ok := c.entries[key]
	return ok
}

func (c TimeCache) Add(key string) {
	c.entries[key] = time.Now().Unix()
}

func (c TimeCache) Remove(key string) {
	delete(c.entries, key)
}

func (c TimeCache) Get(key string) (int64, bool) {
	val, ok := c.entries[key]
	return val, ok
}

// Writes the current cache to disk, kept at DataPath()/<name>.cache
func (c *TimeCache) Persist() error {
	logger := log.Default()
	logger.Debug("Persisting cache")

	// A smarter way to do this once the file gets large is to append after add
	f, err := os.Create(path.Join(config.DataPath(), c.name+".cache"))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	data, err := c.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	logger.Info("Persisted cache", "name", c.name, "entries", len(c.entries))
	return nil
}

// Loads cache data from disk, kept at DataPath()/<name>.cache
func (c *TimeCache) Load() error {
	logger := log.Default()

	filepath := path.Join(config.DataPath(), c.name+".cache")
	logger.Info("Loading cache from disk", "name", c.name, "path", filepath)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	if err := c.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	logger.Info("Loaded cache from disk", "name", c.name, "entries", len(c.entries))
	return nil
}

func (c *TimeCache) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)

	// Manually encode each field you want to persist
	if err := enc.Encode(c.name); err != nil {
		return nil, fmt.Errorf("failed to encode name: %w", err)
	}
	if err := enc.Encode(c.entries); err != nil {
		return nil, fmt.Errorf("failed to encode entries: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *TimeCache) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	// Initialize the map before decoding
	c.entries = make(map[string]int64)

	// Manually decode each field in the same order as they were encoded
	if err := dec.Decode(&c.name); err != nil {
		return fmt.Errorf("failed to decode name: %w", err)
	}
	if err := dec.Decode(&c.entries); err != nil {
		return fmt.Errorf("failed to decode entries: %w", err)
	}

	return nil
}
