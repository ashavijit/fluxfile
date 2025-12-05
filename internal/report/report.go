package report

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

type TaskResult struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Duration  time.Duration `json:"duration_ns"`
	CacheHit  bool          `json:"cache_hit"`
	Error     string        `json:"error,omitempty"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
}

type Report struct {
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	TotalTime  time.Duration `json:"total_time_ns"`
	Tasks      []TaskResult  `json:"tasks"`
	TotalTasks int           `json:"total_tasks"`
	Passed     int           `json:"passed"`
	Failed     int           `json:"failed"`
	Cached     int           `json:"cached"`
	Skipped    int           `json:"skipped"`
}

type Collector struct {
	mu        sync.Mutex
	results   []TaskResult
	startTime time.Time
}

func NewCollector() *Collector {
	return &Collector{
		results:   make([]TaskResult, 0),
		startTime: time.Now(),
	}
}

func (c *Collector) Add(name string, duration time.Duration, success bool, cached bool, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := TaskResult{
		Name:      name,
		Duration:  duration,
		CacheHit:  cached,
		StartTime: time.Now().Add(-duration),
		EndTime:   time.Now(),
	}

	if cached {
		result.Status = "cached"
	} else if success {
		result.Status = "passed"
	} else {
		result.Status = "failed"
		if err != nil {
			result.Error = err.Error()
		}
	}

	c.results = append(c.results, result)
}

func (c *Collector) AddSkipped(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.results = append(c.results, TaskResult{
		Name:      name,
		Status:    "skipped",
		StartTime: time.Now(),
		EndTime:   time.Now(),
	})
}

func (c *Collector) Generate() *Report {
	c.mu.Lock()
	defer c.mu.Unlock()

	endTime := time.Now()
	report := &Report{
		StartTime:  c.startTime,
		EndTime:    endTime,
		TotalTime:  endTime.Sub(c.startTime),
		Tasks:      c.results,
		TotalTasks: len(c.results),
	}

	for _, r := range c.results {
		switch r.Status {
		case "passed":
			report.Passed++
		case "failed":
			report.Failed++
		case "cached":
			report.Cached++
		case "skipped":
			report.Skipped++
		}
	}

	return report
}

func (r *Report) WriteJSON(path string) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
