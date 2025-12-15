package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	called := false
	callback := func() {
		called = true
	}

	w, err := New([]string{"*.go"}, callback)
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	if w == nil {
		t.Fatal("Expected watcher, got nil")
	}

	if w.watcher == nil {
		t.Error("Expected fsnotify watcher to be initialized")
	}

	if len(w.patterns) != 1 {
		t.Errorf("Expected 1 pattern, got %d", len(w.patterns))
	}

	if w.debounce != 100*time.Millisecond {
		t.Errorf("Expected debounce 100ms, got %v", w.debounce)
	}

	// Callback should not be called yet
	if called {
		t.Error("Callback should not be called during initialization")
	}
}

func TestStop(t *testing.T) {
	w, err := New([]string{"*.txt"}, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}

	if err := w.Stop(); err != nil {
		t.Errorf("Stop failed: %v", err)
	}
}

func TestExpandPatterns(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	goFile1 := filepath.Join(dir, "main.go")
	goFile2 := filepath.Join(dir, "util.go")
	txtFile := filepath.Join(dir, "readme.txt")

	os.WriteFile(goFile1, []byte("package main"), 0644)
	os.WriteFile(goFile2, []byte("package util"), 0644)
	os.WriteFile(txtFile, []byte("readme"), 0644)

	w, err := New([]string{filepath.Join(dir, "*.go")}, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	files, err := w.expandPatterns()
	if err != nil {
		t.Fatalf("expandPatterns failed: %v", err)
	}

	// Should find 2 .go files
	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Verify they're absolute paths
	for _, f := range files {
		if !filepath.IsAbs(f) {
			t.Errorf("Expected absolute path, got %s", f)
		}
	}
}

func TestExpandPatternsMultiple(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("key: value"), 0644)

	patterns := []string{
		filepath.Join(dir, "*.go"),
		filepath.Join(dir, "*.yaml"),
	}

	w, err := New(patterns, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	files, _ := w.expandPatterns()
	if len(files) != 2 {
		t.Errorf("Expected 2 files from multiple patterns, got %d", len(files))
	}
}

func TestExpandPatternsDeduplication(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "main.go")
	os.WriteFile(file, []byte("package main"), 0644)

	patterns := []string{
		filepath.Join(dir, "*.go"),
		filepath.Join(dir, "*.go"),
	}

	w, err := New(patterns, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	files, _ := w.expandPatterns()
	if len(files) != 1 {
		t.Errorf("Expected 1 file after deduplication, got %d", len(files))
	}
}

func TestExpandPatternsNonExistent(t *testing.T) {
	w, err := New([]string{"/nonexistent/path/*.go"}, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	files, err := w.expandPatterns()
	if err != nil {
		t.Fatalf("expandPatterns failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files for non-existent path, got %d", len(files))
	}
}

func TestWatcherWithLogger(t *testing.T) {
	w, err := New([]string{"*.go"}, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	if w.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestWatcherDebounce(t *testing.T) {
	w, err := New([]string{"*.go"}, func() {})
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}
	defer w.Stop()

	if w.debounce != 100*time.Millisecond {
		t.Errorf("Expected 100ms debounce, got %v", w.debounce)
	}
}
