package lock

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
)

type LockFile struct {
	Version   string              `json:"version"`
	Generated time.Time           `json:"generated"`
	Tasks     map[string]TaskLock `json:"tasks"`
}

type TaskLock struct {
	Inputs  map[string]string `json:"inputs"`
	Outputs map[string]string `json:"outputs"`
	Hash    string            `json:"hash"`
}

func Generate(fluxFile *ast.FluxFile) (*LockFile, error) {
	lock := &LockFile{
		Version:   "1.0",
		Generated: time.Now(),
		Tasks:     make(map[string]TaskLock),
	}

	for _, task := range fluxFile.Tasks {
		if len(task.Inputs) == 0 && len(task.Outputs) == 0 {
			continue
		}

		taskLock := TaskLock{
			Inputs:  make(map[string]string),
			Outputs: make(map[string]string),
		}

		for _, pattern := range task.Inputs {
			files, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			for _, file := range files {
				hash, err := hashFile(file)
				if err == nil {
					taskLock.Inputs[file] = hash
				}
			}
		}

		for _, pattern := range task.Outputs {
			files, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			for _, file := range files {
				hash, err := hashFile(file)
				if err == nil {
					taskLock.Outputs[file] = hash
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

		for file, expectedHash := range taskLock.Inputs {
			actualHash, err := hashFile(file)
			if err != nil {
				taskChanges = append(taskChanges, fmt.Sprintf("input %s: missing or unreadable", file))
				continue
			}
			if actualHash != expectedHash {
				taskChanges = append(taskChanges, fmt.Sprintf("input %s: hash mismatch", file))
			}
		}

		for file, expectedHash := range taskLock.Outputs {
			actualHash, err := hashFile(file)
			if err != nil {
				taskChanges = append(taskChanges, fmt.Sprintf("output %s: missing", file))
				continue
			}
			if actualHash != expectedHash {
				taskChanges = append(taskChanges, fmt.Sprintf("output %s: hash mismatch", file))
			}
		}

		if len(taskChanges) > 0 {
			changes[taskName] = taskChanges
		}
	}

	return changes, nil
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
