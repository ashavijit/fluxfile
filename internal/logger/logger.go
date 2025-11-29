package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	verbose bool
}

func New() *Logger {
	return &Logger{verbose: true}
}

func (l *Logger) SetVerbose(v bool) {
	l.verbose = v
}

func (l *Logger) Info(msg string) {
	fmt.Printf("[\033[34mINFO\033[0m] %s\n", msg)
}

func (l *Logger) Warn(msg string) {
	fmt.Printf("[\033[33mWARN\033[0m] %s\n", msg)
}

func (l *Logger) Error(msg string) {
	fmt.Fprintf(os.Stderr, "[\033[31mERROR\033[0m] %s\n", msg)
}

func (l *Logger) TaskStart(name string) {
	fmt.Printf("[\033[36m→\033[0m] Running task: \033[1m%s\033[0m\n", name)
}

func (l *Logger) TaskComplete(name string, duration time.Duration) {
	fmt.Printf("[\033[32m✓\033[0m] Task \033[1m%s\033[0m completed in %v\n", name, duration)
}

func (l *Logger) TaskFailed(name string, err error) {
	fmt.Fprintf(os.Stderr, "[\033[31m✗\033[0m] Task \033[1m%s\033[0m failed: %v\n", name, err)
}

func (l *Logger) TaskCached(name string) {
	fmt.Printf("[\033[35m⚡\033[0m] Task \033[1m%s\033[0m (cached)\n", name)
}

func (l *Logger) Command(cmd string) {
	if l.verbose {
		fmt.Printf("  \033[90m$ %s\033[0m\n", cmd)
	}
}

func (l *Logger) Stdout(line string) {
	fmt.Println("  " + line)
}

func (l *Logger) Stderr(line string) {
	fmt.Fprintln(os.Stderr, "  "+line)
}

func (l *Logger) Fatal(msg string) {
	l.Error(msg)
	os.Exit(1)
}
