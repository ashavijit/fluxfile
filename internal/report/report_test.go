package report

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCollector(t *testing.T) {
	c := NewCollector()

	c.Add("build", 100*time.Millisecond, true, false, nil)
	c.Add("test", 200*time.Millisecond, true, false, nil)
	c.Add("lint", 50*time.Millisecond, true, true, nil)

	report := c.Generate()

	if report.TotalTasks != 3 {
		t.Errorf("Expected 3 tasks, got %d", report.TotalTasks)
	}

	if report.Passed != 2 {
		t.Errorf("Expected 2 passed, got %d", report.Passed)
	}

	if report.Cached != 1 {
		t.Errorf("Expected 1 cached, got %d", report.Cached)
	}
}

func TestCollectorWithFailure(t *testing.T) {
	c := NewCollector()

	c.Add("build", 100*time.Millisecond, false, false, nil)

	report := c.Generate()

	if report.Failed != 1 {
		t.Errorf("Expected 1 failed, got %d", report.Failed)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		input    time.Duration
		contains string
	}{
		{500 * time.Microsecond, "us"},
		{500 * time.Millisecond, "ms"},
		{2 * time.Second, "s"},
		{2 * time.Minute, "m"},
	}

	for _, tt := range tests {
		got := FormatDuration(tt.input)
		if got == "" {
			t.Errorf("FormatDuration(%v) returned empty string", tt.input)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	c := NewCollector()
	c.Add("build", 100*time.Millisecond, true, false, nil)

	report := c.Generate()

	dir := t.TempDir()
	path := filepath.Join(dir, "report.json")

	if err := report.WriteJSON(path); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("JSON file is empty")
	}
}

func TestFormatSummary(t *testing.T) {
	c := NewCollector()
	c.Add("build", 100*time.Millisecond, true, false, nil)
	c.Add("test", 50*time.Millisecond, true, true, nil)

	report := c.Generate()
	summary := FormatSummary(report)

	if summary == "" {
		t.Error("Summary is empty")
	}
}
