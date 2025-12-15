package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/cache"
)

func (e *Executor) evaluateCondition(condition string, vars map[string]string) (bool, error) {
	if condition == "" {
		return true, nil
	}

	condition = strings.TrimSpace(condition)
	condition = strings.ReplaceAll(condition, " = = ", "==")
	condition = strings.ReplaceAll(condition, " ! = ", "!=")
	condition = strings.ReplaceAll(condition, " > = ", ">=")
	condition = strings.ReplaceAll(condition, " < = ", "<=")

	operators := []string{"==", "!=", ">=", "<=", ">", "<"}
	var operator string
	var leftPart, rightPart string

	for _, op := range operators {
		if strings.Contains(condition, op) {
			parts := strings.SplitN(condition, op, 2)
			if len(parts) == 2 {
				operator = op
				leftPart = strings.TrimSpace(parts[0])
				rightPart = strings.TrimSpace(parts[1])
				break
			}
		}
	}

	if operator == "" {
		return false, fmt.Errorf("no valid operator found in condition: %s", condition)
	}

	left := vars[leftPart]
	if left == "" {
		left = leftPart
	}

	right := strings.Trim(rightPart, "\"'")

	switch operator {
	case "==":
		return left == right, nil
	case "!=":
		return left != right, nil
	case ">":
		leftNum, err1 := strconv.Atoi(left)
		rightNum, err2 := strconv.Atoi(right)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("invalid numeric comparison")
		}
		return leftNum > rightNum, nil
	case "<":
		leftNum, err1 := strconv.Atoi(left)
		rightNum, err2 := strconv.Atoi(right)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("invalid numeric comparison")
		}
		return leftNum < rightNum, nil
	case ">=":
		leftNum, err1 := strconv.Atoi(left)
		rightNum, err2 := strconv.Atoi(right)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("invalid numeric comparison")
		}
		return leftNum >= rightNum, nil
	case "<=":
		leftNum, err1 := strconv.Atoi(left)
		rightNum, err2 := strconv.Atoi(right)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("invalid numeric comparison")
		}
		return leftNum <= rightNum, nil
	default:
		return false, fmt.Errorf("unknown operator: %s", operator)
	}
}

func (e *Executor) checkPreconditions(preconditions []ast.Precondition) error {
	for _, pre := range preconditions {
		switch pre.Type {
		case "file":
			if _, err := os.Stat(pre.Value); os.IsNotExist(err) {
				return fmt.Errorf("precondition failed: file %s does not exist", pre.Value)
			}
		case "command":
			if _, err := execLookPath(pre.Value); err != nil {
				return fmt.Errorf("precondition failed: command %s not found", pre.Value)
			}
		case "env":
			if os.Getenv(pre.Value) == "" {
				return fmt.Errorf("precondition failed: environment variable %s not set", pre.Value)
			}
		default:
			return fmt.Errorf("unknown precondition type: %s", pre.Type)
		}
	}
	return nil
}

func (e *Executor) executeWithRetry(task *ast.Task, vars map[string]string) error {
	maxRetries := task.Retries
	if maxRetries <= 0 {
		maxRetries = 1
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := parseRetryDelay(task.RetryDelay)
			e.logger.Info(fmt.Sprintf("Retry attempt %d/%d after %v", attempt+1, maxRetries, delay))
			time.Sleep(delay)
		}

		err := e.runCommands(task, vars)
		if err == nil {
			return nil
		}
		lastErr = err
	}

	return lastErr
}

func (e *Executor) executeWithTimeout(task *ast.Task, vars map[string]string) error {
	if task.Timeout == "" {
		return e.executeWithRetry(task, vars)
	}

	timeout, err := time.ParseDuration(task.Timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout: %s", task.Timeout)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- e.executeWithRetry(task, vars)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("task timed out after %s", task.Timeout)
	}
}

func (e *Executor) runCommands(task *ast.Task, vars map[string]string) error {
	expandedRun := expandSlice(task.Run, vars)
	for _, cmd := range expandedRun {
		if err := e.runCommand(cmd, vars); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) checkEnhancedCache(task *ast.Task, useCache bool) (bool, string) {
	if !useCache || !task.Cache {
		return false, ""
	}

	if len(task.Inputs) == 0 {
		return false, ""
	}

	inputHash, err := cache.HashFiles(task.Inputs)
	if err != nil {
		return false, ""
	}

	if entry, ok := e.cache.Get(task.Name, inputHash); ok && entry.Success {
		if len(task.Outputs) > 0 {
			for _, output := range task.Outputs {
				if _, err := os.Stat(output); os.IsNotExist(err) {
					return false, ""
				}
			}
		}
		return true, inputHash
	}

	return false, inputHash
}

func parseRetryDelay(delayStr string) time.Duration {
	if delayStr == "" {
		return 1 * time.Second
	}
	duration, err := time.ParseDuration(delayStr)
	if err != nil {
		return 1 * time.Second
	}
	return duration
}

func expandSlice(slice []string, vars map[string]string) []string {
	result := make([]string, len(slice))
	for i, s := range slice {
		result[i] = expandString(s, vars)
	}
	return result
}

func expandString(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "${"+k+"}", v)
		s = strings.ReplaceAll(s, "$"+k, v)
	}
	return s
}

func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}
