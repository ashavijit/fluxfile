package watcher

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/ashavijit/fluxfile/internal/logger"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher  *fsnotify.Watcher
	patterns []string
	callback func()
	logger   *logger.Logger
	debounce time.Duration
}

func New(patterns []string, callback func()) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		watcher:  w,
		patterns: patterns,
		callback: callback,
		logger:   logger.New(),
		debounce: 100 * time.Millisecond,
	}, nil
}

func (w *Watcher) Start() error {
	files, err := w.expandPatterns()
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := w.watcher.Add(file); err != nil {
			w.logger.Warn(fmt.Sprintf("Failed to watch %s: %v", file, err))
		}
	}

	w.logger.Info(fmt.Sprintf("Watching %d files...", len(files)))

	var timer *time.Timer
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Create == fsnotify.Create {

				if timer != nil {
					timer.Stop()
				}

				timer = time.AfterFunc(w.debounce, func() {
					w.logger.Info(fmt.Sprintf("File changed: %s", event.Name))
					w.callback()
				})
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			w.logger.Error(fmt.Sprintf("Watcher error: %v", err))
		}
	}
}

func (w *Watcher) Stop() error {
	return w.watcher.Close()
}

func (w *Watcher) expandPatterns() ([]string, error) {
	var files []string
	seen := make(map[string]bool)

	for _, pattern := range w.patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			absPath, err := filepath.Abs(match)
			if err != nil {
				continue
			}

			if !seen[absPath] {
				seen[absPath] = true
				files = append(files, absPath)
			}
		}
	}

	return files, nil
}
