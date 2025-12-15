package lock

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
)

func TestGenerate(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	buildTask := ast.NewTask("build")
	buildTask.Inputs = []string{"*.go"}
	buildTask.Outputs = []string{"bin/app"}
	fluxFile.Tasks = []ast.Task{buildTask}

	lock, err := Generate(fluxFile, "1.0.0")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if lock == nil {
		t.Fatal("Expected lock file, got nil")
	}

	if lock.Version != "2.0" {
		t.Errorf("Expected version 2.0, got %s", lock.Version)
	}

	if lock.Metadata.FluxVersion != "1.0.0" {
		t.Errorf("Expected flux version 1.0.0, got %s", lock.Metadata.FluxVersion)
	}
}

func TestGenerateWithPath(t *testing.T) {
	fluxFile := ast.NewFluxFile()

	lock, err := GenerateWithPath(fluxFile, "custom/FluxFile", "2.0.0")
	if err != nil {
		t.Fatalf("GenerateWithPath failed: %v", err)
	}

	if lock.Metadata.FluxFilePath != "custom/FluxFile" {
		t.Errorf("Expected custom path, got %s", lock.Metadata.FluxFilePath)
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	lockPath := filepath.Join(dir, "FluxFile.lock")

	original := &LockFile{
		Version:   "2.0",
		Generated: time.Now(),
		Metadata: Metadata{
			FluxFilePath: "FluxFile",
			FluxVersion:  "1.0.0",
			OS:           "linux",
			Arch:         "amd64",
		},
		Tasks: map[string]TaskLock{
			"build": {
				ConfigHash:  "abc123",
				CommandHash: "def456",
				Inputs:      map[string]FileInfo{},
				Outputs:     map[string]FileInfo{},
				Hash:        "hash123",
				LastUpdated: time.Now(),
			},
		},
	}

	// Save
	if err := Save(original, lockPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatal("Lock file was not created")
	}

	// Load
	loaded, err := Load(lockPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Version != original.Version {
		t.Errorf("Version mismatch: expected %s, got %s", original.Version, loaded.Version)
	}

	if loaded.Metadata.FluxVersion != original.Metadata.FluxVersion {
		t.Errorf("FluxVersion mismatch")
	}

	if len(loaded.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(loaded.Tasks))
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, err := Load("/nonexistent/path/FluxFile.lock")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestVerify(t *testing.T) {
	dir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(dir, "test.txt")
	os.WriteFile(testFile, []byte("content"), 0644)

	// Get file info
	info, _ := getFileInfo(testFile)

	lock := &LockFile{
		Version: "2.0",
		Tasks: map[string]TaskLock{
			"build": {
				Inputs: map[string]FileInfo{
					testFile: info,
				},
			},
		},
	}

	// Verify should pass initially
	changes, err := Verify(lock)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if len(changes) != 0 {
		t.Error("Expected no changes for unchanged files")
	}

	// Modify the file
	os.WriteFile(testFile, []byte("modified content"), 0644)

	// Verify should detect changes
	changes, _ = Verify(lock)
	if len(changes) == 0 {
		t.Error("Expected changes for modified file")
	}
}

func TestVerifyMissingFile(t *testing.T) {
	lock := &LockFile{
		Version: "2.0",
		Tasks: map[string]TaskLock{
			"build": {
				Inputs: map[string]FileInfo{
					"/nonexistent/file.txt": {
						Hash: "abc123",
						Size: 100,
					},
				},
			},
		},
	}

	changes, _ := Verify(lock)
	if len(changes) == 0 {
		t.Error("Expected changes for missing file")
	}

	if _, ok := changes["build"]; !ok {
		t.Error("Expected build task to have changes")
	}
}

func TestUpdateTask(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	task := ast.NewTask("build")
	task.Run = []string{"go build"}
	fluxFile.Tasks = []ast.Task{task}

	lock := &LockFile{
		Version: "2.0",
		Tasks:   make(map[string]TaskLock),
	}

	err := UpdateTask(lock, fluxFile, "build")
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	if _, ok := lock.Tasks["build"]; !ok {
		t.Error("Expected build task to be in lock")
	}
}

func TestUpdateTaskNonExistent(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	lock := &LockFile{Tasks: make(map[string]TaskLock)}

	err := UpdateTask(lock, fluxFile, "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestClean(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Tasks = []ast.Task{ast.NewTask("build")}

	lock := &LockFile{
		Version: "2.0",
		Tasks: map[string]TaskLock{
			"build": {}, // exists in FluxFile
			"old":   {}, // doesn't exist in FluxFile
			"stale": {}, // doesn't exist in FluxFile
		},
	}

	removed := Clean(lock, fluxFile)
	if removed != 2 {
		t.Errorf("Expected to remove 2 tasks, removed %d", removed)
	}

	if _, ok := lock.Tasks["build"]; !ok {
		t.Error("Build task should still exist")
	}

	if _, ok := lock.Tasks["old"]; ok {
		t.Error("Old task should be removed")
	}

	if _, ok := lock.Tasks["stale"]; ok {
		t.Error("Stale task should be removed")
	}
}

func TestComputeDiff(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "src.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	info, _ := getFileInfo(testFile)

	lock := &LockFile{
		Version: "2.0",
		Tasks: map[string]TaskLock{
			"build": {
				ConfigHash:  "oldconfig",
				CommandHash: "oldcommand",
				Inputs: map[string]FileInfo{
					testFile: info,
				},
			},
		},
	}

	fluxFile := ast.NewFluxFile()
	task := ast.NewTask("build")
	task.Run = []string{"go build"} // Different command
	fluxFile.Tasks = []ast.Task{task}

	diffs := ComputeDiff(lock, fluxFile)

	// Should detect command change
	if len(diffs) == 0 {
		t.Error("Expected diffs for changed command")
	}

	if len(diffs) > 0 && !diffs[0].CommandChanged {
		t.Error("Expected command change to be detected")
	}
}

func TestHashFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")

	os.WriteFile(file, []byte("hello world"), 0644)

	hash1, err := hashFile(file)
	if err != nil {
		t.Fatalf("hashFile failed: %v", err)
	}

	if hash1 == "" {
		t.Error("Expected non-empty hash")
	}

	// Same content should give same hash
	hash2, _ := hashFile(file)
	if hash1 != hash2 {
		t.Error("Same file should give same hash")
	}

	// Different content should give different hash
	os.WriteFile(file, []byte("different"), 0644)
	hash3, _ := hashFile(file)
	if hash1 == hash3 {
		t.Error("Different content should give different hash")
	}
}

func TestComputeTaskConfigHash(t *testing.T) {
	task1 := ast.NewTask("build")
	task1.Deps = []string{"clean"}
	task1.Parallel = true

	task2 := ast.NewTask("build")
	task2.Deps = []string{"clean"}
	task2.Parallel = true

	task3 := ast.NewTask("build")
	task3.Deps = []string{"different"}
	task3.Parallel = true

	hash1 := computeTaskConfigHash(task1)
	hash2 := computeTaskConfigHash(task2)
	hash3 := computeTaskConfigHash(task3)

	if hash1 != hash2 {
		t.Error("Same config should give same hash")
	}

	if hash1 == hash3 {
		t.Error("Different config should give different hash")
	}
}

func TestComputeCommandHash(t *testing.T) {
	cmd1 := []string{"go build", "go test"}
	cmd2 := []string{"go build", "go test"}
	cmd3 := []string{"npm install", "npm test"}

	hash1 := computeCommandHash(cmd1)
	hash2 := computeCommandHash(cmd2)
	hash3 := computeCommandHash(cmd3)

	if hash1 != hash2 {
		t.Error("Same commands should give same hash")
	}

	if hash1 == hash3 {
		t.Error("Different commands should give different hash")
	}
}

func TestGetFileInfo(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	content := []byte("test content")

	os.WriteFile(file, content, 0644)

	info, err := getFileInfo(file)
	if err != nil {
		t.Fatalf("getFileInfo failed: %v", err)
	}

	if info.Hash == "" {
		t.Error("Expected non-empty hash")
	}

	if info.Size != int64(len(content)) {
		t.Errorf("Expected size %d, got %d", len(content), info.Size)
	}

	if info.ModTime.IsZero() {
		t.Error("Expected non-zero mod time")
	}
}

func TestGetFileInfoNonExistent(t *testing.T) {
	_, err := getFileInfo("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
