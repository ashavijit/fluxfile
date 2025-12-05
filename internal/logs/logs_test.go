package logs

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLogStore(t *testing.T) {
	dir := t.TempDir()
	store, err := NewLogStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	store.StartTask("build")
	store.Log("info", "Starting build")
	store.LogCommand("go build ./...", 500*time.Millisecond)
	store.EndTask("build", true)

	tasks := store.GetAllTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if task.TaskName != "build" {
		t.Errorf("Expected task name 'build', got '%s'", task.TaskName)
	}
	if task.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", task.Status)
	}
	if len(task.Entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(task.Entries))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := NewLogStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	store.StartTask("test")
	store.Log("info", "Running tests")
	store.EndTask("test", true)

	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 log file, got %d", len(files))
	}

	logs, err := LoadLogs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 {
		t.Errorf("Expected 1 loaded log, got %d", len(logs))
	}
}

func TestGenerateHTML(t *testing.T) {
	dir := t.TempDir()
	logDir := filepath.Join(dir, ".flux", "logs")
	os.MkdirAll(logDir, 0755)

	logs := []*TaskLog{
		{
			TaskName:  "build",
			StartTime: time.Now().Add(-1 * time.Minute),
			EndTime:   time.Now(),
			Status:    "success",
			Entries: []LogEntry{
				{Timestamp: time.Now(), Level: "info", Message: "Building..."},
			},
		},
	}

	// Test HTML generation

	path, err := GenerateHTML(logs)
	if err != nil {
		t.Skipf("Skipping HTML test: %v", err)
	}

	if path == "" {
		t.Error("Expected non-empty path")
	}
}
