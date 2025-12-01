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
	Status   string // pending, running, success, failed
	Duration time.Duration
	Error    error
	mu       sync.RWMutex
}

type TUIState struct {
	Tasks  map[string]*TaskStatus
	Output []string
	mu     sync.RWMutex
}

func runInteractiveTUI(exec *executor.Executor, taskName string, profile string, useCache bool) {
	state := &TUIState{
		Tasks:  make(map[string]*TaskStatus),
		Output: []string{},
	}

	clearScreen()
	printHeader()

	tasks := exec.ListTasks()
	for _, name := range tasks {
		state.Tasks[name] = &TaskStatus{
			Name:   name,
			Status: "pending",
		}
	}

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			clearScreen()
			printHeader()
			state.render()
		}
	}()

	start := time.Now()
	err := exec.Execute(taskName, profile, useCache)

	clearScreen()
	printHeader()

	if err != nil {
		fmt.Printf("\n%s✗ Task failed: %s%s\n", colorRed, taskName, colorReset)
		fmt.Printf("%sError: %s%s\n\n", colorRed, err.Error(), colorReset)
	} else {
		fmt.Printf("\n%s✓ Task completed: %s%s\n", colorGreen, taskName, colorReset)
	}

	fmt.Printf("%sDuration: %v%s\n\n", colorGray, time.Since(start), colorReset)
}

func (s *TUIState) render() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                      TASK EXECUTION STATUS                     ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", colorCyan, colorReset)
	fmt.Println()

	maxNameLen := 15
	for _, task := range s.Tasks {
		if len(task.Name) > maxNameLen {
			maxNameLen = len(task.Name)
		}
	}

	fmt.Printf("  %s%-*s  %-12s  %s%s\n", colorYellow, maxNameLen, "TASK", "STATUS", "DURATION", colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, strings.Repeat("─", maxNameLen+30), colorReset)

	for _, task := range s.Tasks {
		task.mu.RLock()
		statusIcon, statusColor := getStatusDisplay(task.Status)
		durationStr := formatDuration(task.Duration)

		fmt.Printf("  %s%-*s%s  %s%-12s%s  %s%s%s\n",
			colorGreen, maxNameLen, task.Name, colorReset,
			statusColor, statusIcon, colorReset,
			colorGray, durationStr, colorReset)
		task.mu.RUnlock()
	}

	fmt.Println()

	if len(s.Output) > 0 {
		fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
		fmt.Printf("%s║                           OUTPUT LOG                           ║%s\n", colorCyan, colorReset)
		fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", colorCyan, colorReset)
		fmt.Println()

		start := 0
		if len(s.Output) > 10 {
			start = len(s.Output) - 10
		}

		for i := start; i < len(s.Output); i++ {
			fmt.Printf("  %s%s%s\n", colorGray, s.Output[i], colorReset)
		}
	}

	fmt.Println()
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

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func printHeader() {
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                       FLUX INTERACTIVE MODE                     ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                     Press Ctrl+C to abort                       ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", colorCyan, colorReset)
}

const (
	colorRed = "\033[31m"
)
