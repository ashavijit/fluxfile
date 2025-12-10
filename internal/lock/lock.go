package lock

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
)

type LockFile struct {
	Version      string              `json:"version"`
	Generated    time.Time           `json:"generated"`
	Metadata     Metadata            `json:"metadata"`
	FluxFileHash string              `json:"fluxfile_hash"`
	Tasks        map[string]TaskLock `json:"tasks"`
}

type Metadata struct {
	FluxFilePath string `json:"fluxfile_path"`
	Hostname     string `json:"hostname,omitempty"`
	User         string `json:"user,omitempty"`
	FluxVersion  string `json:"flux_version"`
	OS           string `json:"os"`
	Arch         string `json:"arch"`
}

type TaskLock struct {
	ConfigHash  string              `json:"config_hash"`
	CommandHash string              `json:"command_hash"`
	Inputs      map[string]FileInfo `json:"inputs"`
	Outputs     map[string]FileInfo `json:"outputs"`
	Hash        string              `json:"hash"`
	LastUpdated time.Time           `json:"last_updated"`
}

type FileInfo struct {
	Hash    string    `json:"hash"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type DiffResult struct {
	TaskName       string
	ConfigChanged  bool
	CommandChanged bool
	InputChanges   []FileChange
	OutputChanges  []FileChange
}

type FileChange struct {
	Path       string
	ChangeType string // "modified", "missing", "new", "size_changed"
	OldHash    string
	NewHash    string
	OldSize    int64
	NewSize    int64
}

func Generate(fluxFile *ast.FluxFile, version string) (*LockFile, error) {
	return GenerateWithPath(fluxFile, "FluxFile", version)
}

func GenerateWithPath(fluxFile *ast.FluxFile, fluxFilePath string, version string) (*LockFile, error) {
	hostname, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME") // Windows
	}

	fluxFileHash, _ := hashFile(fluxFilePath)

	lock := &LockFile{
		Version:      "2.0",
		Generated:    time.Now(),
		FluxFileHash: fluxFileHash,
		Metadata: Metadata{
			FluxFilePath: fluxFilePath,
			Hostname:     hostname,
			User:         user,
			FluxVersion:  version,
			OS:           runtime.GOOS,
			Arch:         runtime.GOARCH,
		},
		Tasks: make(map[string]TaskLock),
	}

	for _, task := range fluxFile.Tasks {
		if len(task.Inputs) == 0 && len(task.Outputs) == 0 {
			continue
		}

		taskLock := TaskLock{
			Inputs:      make(map[string]FileInfo),
			Outputs:     make(map[string]FileInfo),
			LastUpdated: time.Now(),
		}

		taskLock.ConfigHash = computeTaskConfigHash(task)
		taskLock.CommandHash = computeCommandHash(task.Run)

		for _, pattern := range task.Inputs {
			files, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			for _, file := range files {
				info, err := getFileInfo(file)
				if err == nil {
					taskLock.Inputs[file] = info
				}
			}
		}

		for _, pattern := range task.Outputs {
			files, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			for _, file := range files {
				info, err := getFileInfo(file)
				if err == nil {
					taskLock.Outputs[file] = info
				}
			}
		}

		taskLock.Hash = computeTaskHash(taskLock)
		lock.Tasks[task.Name] = taskLock
	}

	return lock, nil
}

func Save(lock *LockFile, path string) error {
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func Load(path string) (*LockFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var lock LockFile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, err
	}

	return &lock, nil
}

func Verify(lock *LockFile) (map[string][]string, error) {
	changes := make(map[string][]string)

	for taskName, taskLock := range lock.Tasks {
		var taskChanges []string

		for file, expectedInfo := range taskLock.Inputs {
			actualInfo, err := getFileInfo(file)
			if err != nil {
				taskChanges = append(taskChanges, fmt.Sprintf("input %s: missing or unreadable", file))
				continue
			}
			if actualInfo.Hash != expectedInfo.Hash {
				taskChanges = append(taskChanges, fmt.Sprintf("input %s: hash mismatch (size: %d -> %d)",
					file, expectedInfo.Size, actualInfo.Size))
			}
		}

		for file, expectedInfo := range taskLock.Outputs {
			actualInfo, err := getFileInfo(file)
			if err != nil {
				taskChanges = append(taskChanges, fmt.Sprintf("output %s: missing", file))
				continue
			}
			if actualInfo.Hash != expectedInfo.Hash {
				taskChanges = append(taskChanges, fmt.Sprintf("output %s: hash mismatch (size: %d -> %d)",
					file, expectedInfo.Size, actualInfo.Size))
			}
		}

		if len(taskChanges) > 0 {
			changes[taskName] = taskChanges
		}
	}

	return changes, nil
}

// ComputeDiff generates a detailed diff between lock file and current state
func ComputeDiff(lock *LockFile, fluxFile *ast.FluxFile) []DiffResult {
	var results []DiffResult

	for _, task := range fluxFile.Tasks {
		taskLock, exists := lock.Tasks[task.Name]
		if !exists {
			continue
		}

		diff := DiffResult{
			TaskName: task.Name,
		}

		// Check config changes
		currentConfigHash := computeTaskConfigHash(task)
		diff.ConfigChanged = currentConfigHash != taskLock.ConfigHash

		// Check command changes
		currentCommandHash := computeCommandHash(task.Run)
		diff.CommandChanged = currentCommandHash != taskLock.CommandHash

		// Check input changes
		for file, oldInfo := range taskLock.Inputs {
			newInfo, err := getFileInfo(file)
			if err != nil {
				diff.InputChanges = append(diff.InputChanges, FileChange{
					Path:       file,
					ChangeType: "missing",
					OldHash:    oldInfo.Hash,
					OldSize:    oldInfo.Size,
				})
			} else if newInfo.Hash != oldInfo.Hash {
				changeType := "modified"
				if newInfo.Size != oldInfo.Size {
					changeType = "size_changed"
				}
				diff.InputChanges = append(diff.InputChanges, FileChange{
					Path:       file,
					ChangeType: changeType,
					OldHash:    oldInfo.Hash,
					NewHash:    newInfo.Hash,
					OldSize:    oldInfo.Size,
					NewSize:    newInfo.Size,
				})
			}
		}

		// Check output changes
		for file, oldInfo := range taskLock.Outputs {
			newInfo, err := getFileInfo(file)
			if err != nil {
				diff.OutputChanges = append(diff.OutputChanges, FileChange{
					Path:       file,
					ChangeType: "missing",
					OldHash:    oldInfo.Hash,
					OldSize:    oldInfo.Size,
				})
			} else if newInfo.Hash != oldInfo.Hash {
				changeType := "modified"
				if newInfo.Size != oldInfo.Size {
					changeType = "size_changed"
				}
				diff.OutputChanges = append(diff.OutputChanges, FileChange{
					Path:       file,
					ChangeType: changeType,
					OldHash:    oldInfo.Hash,
					NewHash:    newInfo.Hash,
					OldSize:    oldInfo.Size,
					NewSize:    newInfo.Size,
				})
			}
		}

		if diff.ConfigChanged || diff.CommandChanged || len(diff.InputChanges) > 0 || len(diff.OutputChanges) > 0 {
			results = append(results, diff)
		}
	}

	return results
}

// UpdateTask updates a specific task in the lock file
func UpdateTask(lock *LockFile, fluxFile *ast.FluxFile, taskName string) error {
	var task *ast.Task
	for i := range fluxFile.Tasks {
		if fluxFile.Tasks[i].Name == taskName {
			task = &fluxFile.Tasks[i]
			break
		}
	}

	if task == nil {
		return fmt.Errorf("task '%s' not found in FluxFile", taskName)
	}

	taskLock := TaskLock{
		Inputs:      make(map[string]FileInfo),
		Outputs:     make(map[string]FileInfo),
		LastUpdated: time.Now(),
	}

	taskLock.ConfigHash = computeTaskConfigHash(*task)
	taskLock.CommandHash = computeCommandHash(task.Run)

	for _, pattern := range task.Inputs {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, file := range files {
			info, err := getFileInfo(file)
			if err == nil {
				taskLock.Inputs[file] = info
			}
		}
	}

	for _, pattern := range task.Outputs {
		files, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, file := range files {
			info, err := getFileInfo(file)
			if err == nil {
				taskLock.Outputs[file] = info
			}
		}
	}

	taskLock.Hash = computeTaskHash(taskLock)
	lock.Tasks[taskName] = taskLock
	lock.Generated = time.Now()

	return nil
}

// Clean removes tasks from lock that are not in the FluxFile
func Clean(lock *LockFile, fluxFile *ast.FluxFile) int {
	taskMap := make(map[string]bool)
	for _, task := range fluxFile.Tasks {
		taskMap[task.Name] = true
	}

	removed := 0
	for taskName := range lock.Tasks {
		if !taskMap[taskName] {
			delete(lock.Tasks, taskName)
			removed++
		}
	}

	if removed > 0 {
		lock.Generated = time.Now()
	}

	return removed
}

func getFileInfo(path string) (FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	hash, err := hashFile(path)
	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		Hash:    hash,
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}, nil
}

func hashFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func computeTaskHash(taskLock TaskLock) string {
	data, _ := json.Marshal(taskLock.Inputs)
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash[:])
}

func computeTaskConfigHash(task ast.Task) string {
	// Create a deterministic representation of task configuration
	var parts []string

	// Add dependencies
	parts = append(parts, fmt.Sprintf("deps:%v", task.Deps))

	// Add environment variables (sorted for consistency)
	if len(task.Env) > 0 {
		var envKeys []string
		for k := range task.Env {
			envKeys = append(envKeys, k)
		}
		sort.Strings(envKeys)
		for _, k := range envKeys {
			parts = append(parts, fmt.Sprintf("env:%s=%s", k, task.Env[k]))
		}
	}

	// Add other config
	parts = append(parts, fmt.Sprintf("parallel:%v", task.Parallel))
	parts = append(parts, fmt.Sprintf("cache:%v", task.Cache))
	parts = append(parts, fmt.Sprintf("docker:%v", task.Docker))
	if task.If != "" {
		parts = append(parts, fmt.Sprintf("if:%s", task.If))
	}
	if task.Remote != "" {
		parts = append(parts, fmt.Sprintf("remote:%s", task.Remote))
	}
	if task.Timeout != "" {
		parts = append(parts, fmt.Sprintf("timeout:%s", task.Timeout))
	}

	data := strings.Join(parts, "|")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}

func computeCommandHash(commands []string) string {
	data := strings.Join(commands, "\n")
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}
