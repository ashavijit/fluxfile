package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	dir := t.TempDir()
	c, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	if c == nil {
		t.Fatal("Expected cache, got nil")
	}

	if c.dir != dir {
		t.Errorf("Expected dir %s, got %s", dir, c.dir)
	}
}

func TestNewCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "cache", "dir")
	c, err := New(dir)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	if _, err := os.Stat(c.dir); os.IsNotExist(err) {
		t.Error("Expected cache directory to be created")
	}
}

func TestSetAndGet(t *testing.T) {
	c, _ := New(t.TempDir())

	entry := &CacheEntry{
		TaskName:   "build",
		InputHash:  "abc123",
		OutputHash: "def456",
		Timestamp:  time.Now(),
		Success:    true,
		Duration:   2 * time.Second,
	}

	// Set entry
	if err := c.Set(entry); err != nil {
		t.Fatalf("Failed to set cache entry: %v", err)
	}

	// Get entry with matching hash
	retrieved, ok := c.Get("build", "abc123")
	if !ok {
		t.Fatal("Expected to find cache entry")
	}

	if retrieved.TaskName != entry.TaskName {
		t.Errorf("Expected TaskName %s, got %s", entry.TaskName, retrieved.TaskName)
	}

	if retrieved.InputHash != entry.InputHash {
		t.Errorf("Expected InputHash %s, got %s", entry.InputHash, retrieved.InputHash)
	}

	if retrieved.Success != entry.Success {
		t.Errorf("Expected Success %v, got %v", entry.Success, retrieved.Success)
	}
}

func TestGetWithWrongHash(t *testing.T) {
	c, _ := New(t.TempDir())

	entry := &CacheEntry{
		TaskName:   "build",
		InputHash:  "abc123",
		OutputHash: "def456",
		Timestamp:  time.Now(),
		Success:    true,
	}

	c.Set(entry)

	// Try to get with different hash
	_, ok := c.Get("build", "wronghash")
	if ok {
		t.Error("Should not find cache entry with wrong hash")
	}
}

func TestGetNonExistent(t *testing.T) {
	c, _ := New(t.TempDir())

	_, ok := c.Get("nonexistent", "somehash")
	if ok {
		t.Error("Should not find non-existent cache entry")
	}
}

func TestClear(t *testing.T) {
	dir := t.TempDir()
	c, _ := New(dir)

	// Set some entries
	c.Set(&CacheEntry{TaskName: "task1", InputHash: "hash1", Success: true})
	c.Set(&CacheEntry{TaskName: "task2", InputHash: "hash2", Success: true})

	// Clear cache
	if err := c.Clear(); err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify directory is removed
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Error("Expected cache directory to be removed")
	}
}

func TestEntryPath(t *testing.T) {
	c, _ := New(t.TempDir())

	path := c.entryPath("build")
	expected := filepath.Join(c.dir, "build.json")

	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestHashString(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"hello"},
		{"world"},
		{""},
		{"go build -o app"},
	}

	seen := make(map[string]bool)

	for _, tt := range tests {
		hash := HashString(tt.input)

		// Should be 64 hex characters (SHA256)
		if len(hash) != 64 {
			t.Errorf("Expected 64 char hash, got %d for input %q", len(hash), tt.input)
		}

		// Same input should give same hash
		hash2 := HashString(tt.input)
		if hash != hash2 {
			t.Errorf("Same input should give same hash")
		}

		// Different inputs should give different hashes
		if tt.input != "" && seen[hash] {
			t.Errorf("Hash collision detected for %q", tt.input)
		}
		seen[hash] = true
	}
}

func TestHashFiles(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(dir, "file2.txt")

	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	// Hash single file
	hash1, err := HashFiles([]string{file1})
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}
	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}

	// Hash multiple files
	hash2, err := HashFiles([]string{file1, file2})
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	// Different file sets should give different hashes
	if hash1 == hash2 {
		t.Error("Different file sets should give different hashes")
	}
}

func TestHashFilesWithGlob(t *testing.T) {
	dir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(dir, "a.go"), []byte("package a"), 0644)
	os.WriteFile(filepath.Join(dir, "b.go"), []byte("package b"), 0644)
	os.WriteFile(filepath.Join(dir, "c.txt"), []byte("text file"), 0644)

	// Hash Go files with glob
	pattern := filepath.Join(dir, "*.go")
	hash, err := HashFiles([]string{pattern})
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash for glob pattern")
	}
}

func TestHashFilesNonExistent(t *testing.T) {
	// Hash non-existent files should return empty hash (gracefully fail)
	hash, err := HashFiles([]string{"/nonexistent/path/*.go"})
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	// Should return a hash (of empty content)
	if hash == "" {
		t.Error("Expected hash even with no matching files")
	}
}

func TestHashFilesIgnoresDirectories(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "subdir")
	os.MkdirAll(subDir, 0755)

	os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)

	// Should not fail when pattern includes directories
	pattern := filepath.Join(dir, "*")
	hash, err := HashFiles([]string{pattern})
	if err != nil {
		t.Fatalf("HashFiles failed: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}
}

func TestCacheEntryTimestamp(t *testing.T) {
	c, _ := New(t.TempDir())
	now := time.Now()

	entry := &CacheEntry{
		TaskName:  "task",
		InputHash: "hash",
		Timestamp: now,
		Success:   true,
		Duration:  5 * time.Second,
	}

	c.Set(entry)
	retrieved, ok := c.Get("task", "hash")
	if !ok {
		t.Fatal("Failed to retrieve entry")
	}

	// Timestamp should be preserved (within tolerance for JSON marshaling)
	diff := retrieved.Timestamp.Sub(now)
	if diff < -time.Second || diff > time.Second {
		t.Errorf("Timestamp not preserved correctly, diff: %v", diff)
	}

	if retrieved.Duration != 5*time.Second {
		t.Errorf("Expected duration 5s, got %v", retrieved.Duration)
	}
}
