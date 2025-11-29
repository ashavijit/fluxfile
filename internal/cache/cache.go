package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Cache struct {
	dir string
}

type CacheEntry struct {
	TaskName   string
	InputHash  string
	OutputHash string
	Timestamp  time.Time
	Success    bool
	Duration   time.Duration
}

func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &Cache{dir: dir}, nil
}

func (c *Cache) Get(taskName string, inputHash string) (*CacheEntry, bool) {
	path := c.entryPath(taskName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, false
	}

	if entry.InputHash != inputHash {
		return nil, false
	}

	return &entry, true
}

func (c *Cache) Set(entry *CacheEntry) error {
	path := c.entryPath(entry.TaskName)
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Cache) entryPath(taskName string) string {
	return filepath.Join(c.dir, fmt.Sprintf("%s.json", taskName))
}

func (c *Cache) Clear() error {
	return os.RemoveAll(c.dir)
}

func HashFiles(patterns []string) (string, error) {
	h := sha256.New()

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			info, err := os.Stat(match)
			if err != nil {
				continue
			}

			if info.IsDir() {
				continue
			}

			f, err := os.Open(match)
			if err != nil {
				continue
			}

			if _, err := io.Copy(h, f); err != nil {
				f.Close()
				continue
			}
			f.Close()

			h.Write([]byte(match))
			h.Write([]byte(info.ModTime().String()))
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func HashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
