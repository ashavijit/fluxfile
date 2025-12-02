package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ashavijit/fluxfile/internal/executor"
)

type TaskStatus struct {
	Name     string
	Status   string
	Duration time.Duration
	Error    error
	mu       sync.RWMutex
}

type TUIState struct {
	Tasks map[string]*TaskStatus
	mu    sync.RWMutex
}

func runInteractiveTUI(exec *executor.Executor, taskName string, profile string, useCache bool) {
	fmt.Print("\033[?25l") // Hide cursor
	defer fmt.Print("\033[?25h")

	state := &TUIState{
		Tasks: make(map[string]*TaskStatus),
	}

	tasks := exec.ListTasks()
	for _, name := range tasks {
		state.Tasks[name] = &TaskStatus{Name: name, Status: "pending"}
	}

	fmt.Print("\033[2J\033[H") // Clear screen

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
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

	start := time.Now()
	if t, ok := state.Tasks[taskName]; ok {
		t.mu.Lock()
		t.Status = "running"
		t.mu.Unlock()
	}

	err := exec.Execute(taskName, profile, useCache)

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
	state.render()
	fmt.Printf("\033[%dB", len(state.Tasks)+8) // Move cursor to bottom
}

func (s *TUIState) render() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Print("\033[H") // Move to top-left

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                    FLUX INTERACTIVE MODE                       ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	maxLen := 15
	sorted := make([]string, 0, len(s.Tasks))
	for name := range s.Tasks {
		if len(name) > maxLen {
			maxLen = len(name)
		}
		sorted = append(sorted, name)
	}

	fmt.Printf("  %s%-*s  %-10s  %s%s\n", colorYellow, maxLen, "TASK", "STATUS", "TIME", colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, strings.Repeat("─", maxLen+25), colorReset)

	for _, name := range sorted {
		t := s.Tasks[name]
		t.mu.RLock()

		icon, col := "⏸", colorGray
		if t.Status == "running" {
			icon, col = "🔄", colorYellow
		}
		if t.Status == "success" {
			icon, col = "✓", colorGreen
		}
		if t.Status == "failed" {
			icon, col = "✗", colorRed
		}

		dur := "-"
		if t.Duration > 0 {
			dur = fmt.Sprintf("%.1fs", t.Duration.Seconds())
		}

		fmt.Printf("  %s%-*s%s  %s%-10s%s  %s%s%s\033[K\n",
			colorGreen, maxLen, t.Name, colorReset,
			col, icon+" "+t.Status, colorReset,
			colorGray, dur, colorReset)
		t.mu.RUnlock()
	}
}

const (
	colorMagenta = "\033[35m"
)
