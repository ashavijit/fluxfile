package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ashavijit/fluxfile/internal/executor"
)

// TaskStatus tracks the state of a task in the TUI
type TaskStatus struct {
	Name     string
	Status   string // pending, running, success, failed
	Duration time.Duration
	Error    error
	mu       sync.RWMutex
}

// TUIState holds the global state of the TUI
type TUIState struct {
	Tasks  map[string]*TaskStatus
	Output []string
	mu     sync.RWMutex
}

func runInteractiveTUI(exec *executor.Executor, taskName string, profile string, useCache bool) {
	// Hide cursor
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit

	state := &TUIState{
		Tasks:  make(map[string]*TaskStatus),
		Output: []string{},
	}

	// Initialize tasks
	tasks := exec.ListTasks()
	for _, name := range tasks {
		state.Tasks[name] = &TaskStatus{
			Name:   name,
			Status: "pending",
		}
	}

	// Clear screen once at start
	fmt.Print("\033[2J\033[H")

	// Start update loop
	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				state.render()
			}
		}
	}()

	// Execute task
	// Note: In a real TUI, we'd want to capture exec output and feed it to state.Output
	// For now, we'll let it run. If exec prints to stdout, it will mess up the TUI.
	// We assume exec is silent or we need to modify exec to be silent/capture output.
	start := time.Now()
	
	// Update status to running for the main task
	if t, ok := state.Tasks[taskName]; ok {
		t.mu.Lock()
		t.Status = "running"
		t.mu.Unlock()
	}

	err := exec.Execute(taskName, profile, useCache)
	
	// Update final status
	if t, ok := state.Tasks[taskName]; ok {
		t.mu.Lock()
		if err != nil {
			t.Status = "failed"
			t.Error = err
		} else {
			t.Status = "success"
		}
		t.Duration = time.Since(start)
		t.mu.Unlock()
	}

	done <- true
	state.render() // Final render
	
	// Move cursor below the table
	fmt.Printf("\033[%dB", len(state.Tasks)+10)
}

func (s *TUIState) render() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Move to top-left
	fmt.Print("\033[H")

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                    FLUX INTERACTIVE MODE                       ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	maxNameLen := 15
	sortedNames := make([]string, 0, len(s.Tasks))
	for name := range s.Tasks {
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}
		sortedNames = append(sortedNames, name)
	}

	fmt.Printf("  %s%-*s  %-12s  %s%s\n", colorYellow, maxNameLen, "TASK", "STATUS", "DURATION", colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, strings.Repeat("─", maxNameLen+30), colorReset)

	for _, name := range sortedNames {
		task := s.Tasks[name]
		task.mu.RLock()
		statusIcon, statusColor := getStatusDisplay(task.Status)
		durationStr := formatDuration(task.Duration)
		
		fmt.Printf("  %s%-*s%s  %s%-12s%s  %s%s%s\033[K\n",
			colorGreen, maxNameLen, task.Name, colorReset,
			statusColor, statusIcon, colorReset,
			colorGray, durationStr, colorReset)
		task.mu.RUnlock()
	}
	
	// Clear remaining lines if any (optional, for cleaner look)
	fmt.Print("\033[J")
}

func getStatusDisplay(status string) (string, string) {
	switch status {
	case "running":
		return "🔄 Running", colorYellow
	case "success":
		return "✓ Success", colorGreen
	case "failed":
		return "✗ Failed", colorRed
	default:
		return "⏸ Pending", colorGray
	}
}

func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorReset  = "\033[0m"
)
