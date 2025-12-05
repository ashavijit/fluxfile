package report

import (
	"fmt"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fus", float64(d.Microseconds()))
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Milliseconds()))
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return fmt.Sprintf("%.2fm", d.Minutes())
}

func (r *Report) Print() {
	fmt.Println()
	fmt.Printf("%s%s Execution Report %s\n", colorBold, colorCyan, colorReset)
	fmt.Println(strings.Repeat("-", 60))

	fmt.Printf("\n%-30s %-10s %-15s\n", "TASK", "STATUS", "DURATION")
	fmt.Println(strings.Repeat("-", 60))

	for _, task := range r.Tasks {
		statusColor := colorGreen
		switch task.Status {
		case "failed":
			statusColor = colorRed
		case "cached":
			statusColor = colorYellow
		case "skipped":
			statusColor = colorGray
		}

		duration := FormatDuration(task.Duration)
		if task.Status == "cached" || task.Status == "skipped" {
			duration = "-"
		}

		fmt.Printf("%-30s %s%-10s%s %-15s\n",
			truncate(task.Name, 30),
			statusColor,
			task.Status,
			colorReset,
			duration,
		)

		if task.Error != "" {
			fmt.Printf("  %s%s%s\n", colorRed, task.Error, colorReset)
		}
	}

	// Summary
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("\n%sSummary:%s\n", colorBold, colorReset)
	fmt.Printf("  Total:   %d tasks\n", r.TotalTasks)
	fmt.Printf("  %sPassed:%s  %d\n", colorGreen, colorReset, r.Passed)
	if r.Failed > 0 {
		fmt.Printf("  %sFailed:%s  %d\n", colorRed, colorReset, r.Failed)
	}
	if r.Cached > 0 {
		fmt.Printf("  %sCached:%s  %d\n", colorYellow, colorReset, r.Cached)
	}
	if r.Skipped > 0 {
		fmt.Printf("  %sSkipped:%s %d\n", colorGray, colorReset, r.Skipped)
	}
	fmt.Printf("\n  Total time: %s%s%s\n", colorCyan, FormatDuration(r.TotalTime), colorReset)
	fmt.Println()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func FormatSummary(r *Report) string {
	return fmt.Sprintf("%d passed, %d failed, %d cached in %s",
		r.Passed, r.Failed, r.Cached, FormatDuration(r.TotalTime))
}
