package main

import (
	"fmt"
	"time"

	"github.com/ashavijit/fluxfile/internal/executor"
)

func runInteractiveTUI(exec *executor.Executor, taskName string, profile string, useCache bool) {
	fmt.Print("\033[2J\033[H") // Clear screen once

	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                    FLUX INTERACTIVE MODE                       ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n\n", colorCyan, colorReset)

	fmt.Printf("  %sTask:%s %s\n", colorGray, colorReset, taskName)
	if profile != "" {
		fmt.Printf("  %sProfile:%s %s\n", colorGray, colorReset, profile)
	}
	fmt.Printf("  %sCache:%s %v\n\n", colorGray, colorReset, useCache)

	fmt.Printf("%s▶ Starting execution...%s\n\n", colorYellow, colorReset)

	start := time.Now()
	err := exec.Execute(taskName, profile, useCache)
	duration := time.Since(start)

	fmt.Println()
	fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n", colorCyan, colorReset)

	if err != nil {
		fmt.Printf("  %s✗ FAILED%s in %v\n", colorRed, colorReset, duration)
		fmt.Printf("  %sError: %s%s\n", colorRed, err.Error(), colorReset)
	} else {
		fmt.Printf("  %s✓ SUCCESS%s in %v\n", colorGreen, colorReset, duration)
	}

	fmt.Printf("%s═══════════════════════════════════════════════════════════════%s\n\n", colorCyan, colorReset)
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
