package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Task      string    `json:"task"`
	Message   string    `json:"message"`
	Command   string    `json:"command,omitempty"`
	Duration  int64     `json:"duration_ms,omitempty"`
	ExitCode  int       `json:"exit_code,omitempty"`
	Output    string    `json:"output,omitempty"`
}

type TaskLog struct {
	TaskName  string     `json:"task_name"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time,omitempty"`
	Status    string     `json:"status"`
	WorkDir   string     `json:"work_dir,omitempty"`
	Profile   string     `json:"profile,omitempty"`
	CacheHit  bool       `json:"cache_hit,omitempty"`
	DepsCount int        `json:"deps_count,omitempty"`
	Error     string     `json:"error,omitempty"`
	Entries   []LogEntry `json:"entries"`
}

type LogStore struct {
	mu      sync.RWMutex
	dir     string
	tasks   map[string]*TaskLog
	current string
}

func NewLogStore(dir string) (*LogStore, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &LogStore{
		dir:   dir,
		tasks: make(map[string]*TaskLog),
	}, nil
}

func (s *LogStore) StartTask(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.current = name
	s.tasks[name] = &TaskLog{
		TaskName:  name,
		StartTime: time.Now(),
		Status:    "running",
		Entries:   make([]LogEntry, 0),
	}
}

func (s *LogStore) Log(level, message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task == nil {
		return
	}

	task.Entries = append(task.Entries, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Task:      s.current,
		Message:   message,
	})
}

func (s *LogStore) LogCommand(command string, duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task == nil {
		return
	}

	task.Entries = append(task.Entries, LogEntry{
		Timestamp: time.Now(),
		Level:     "cmd",
		Task:      s.current,
		Command:   command,
		Duration:  duration.Milliseconds(),
	})
}

func (s *LogStore) EndTask(name string, success bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task := s.tasks[name]
	if task == nil {
		return
	}

	task.EndTime = time.Now()
	if success {
		task.Status = "success"
	} else {
		task.Status = "failed"
	}
}

func (s *LogStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, task := range s.tasks {
		data, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%s_%s.json",
			name,
			task.StartTime.Format("20060102_150405"))
		path := filepath.Join(s.dir, filename)

		if err := os.WriteFile(path, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func (s *LogStore) GetAllTasks() []*TaskLog {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*TaskLog, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}

func LoadLogs(dir string) ([]*TaskLog, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return nil, err
	}

	logs := make([]*TaskLog, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var log TaskLog
		if err := json.Unmarshal(data, &log); err != nil {
			continue
		}
		logs = append(logs, &log)
	}

	return logs, nil
}

func GetLogDir() string {
	return filepath.Join(".flux", "logs")
}

func ClearLogs() (int, error) {
	dir := GetLogDir()
	files, err := filepath.Glob(filepath.Join(dir, "*.json"))
	if err != nil {
		return 0, err
	}

	count := 0
	for _, file := range files {
		if err := os.Remove(file); err == nil {
			count++
		}
	}

	// Also remove HTML file
	htmlPath := filepath.Join(dir, "logs.html")
	os.Remove(htmlPath)

	return count, nil
}

func (s *LogStore) SetTaskInfo(workDir, profile string, depsCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task == nil {
		return
	}

	task.WorkDir = workDir
	task.Profile = profile
	task.DepsCount = depsCount
}

func (s *LogStore) SetCacheHit(hit bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task != nil {
		task.CacheHit = hit
	}
}

func (s *LogStore) SetError(err string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task != nil {
		task.Error = err
	}
}

func (s *LogStore) LogCommandWithOutput(command string, duration time.Duration, exitCode int, output string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current == "" {
		return
	}

	task := s.tasks[s.current]
	if task == nil {
		return
	}

	task.Entries = append(task.Entries, LogEntry{
		Timestamp: time.Now(),
		Level:     "cmd",
		Task:      s.current,
		Command:   command,
		Duration:  duration.Milliseconds(),
		ExitCode:  exitCode,
		Output:    output,
	})
}
