package executor

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/cache"
	"github.com/ashavijit/fluxfile/internal/graph"
	"github.com/ashavijit/fluxfile/internal/logger"
	"github.com/ashavijit/fluxfile/internal/logs"
	"github.com/ashavijit/fluxfile/internal/report"
	"github.com/ashavijit/fluxfile/internal/vars"
)

type Executor struct {
	fluxFile  *ast.FluxFile
	graph     *graph.Graph
	cache     *cache.Cache
	logger    *logger.Logger
	vars      map[string]string
	dryRun    bool
	collector *report.Collector
	logStore  *logs.LogStore
}

type ExecutionResult struct {
	TaskName string
	Success  bool
	Duration time.Duration
	Error    error
}

func New(fluxFile *ast.FluxFile, cacheDir string, dryRun bool) (*Executor, error) {
	g, err := graph.BuildGraph(fluxFile.Tasks)
	if err != nil {
		return nil, err
	}

	c, err := cache.New(cacheDir)
	if err != nil {
		return nil, err
	}

	return &Executor{
		fluxFile: fluxFile,
		graph:    g,
		cache:    c,
		logger:   logger.New(),
		vars:     fluxFile.Vars,
		dryRun:   dryRun,
	}, nil
}

func (e *Executor) SetCollector(c *report.Collector) {
	e.collector = c
}

func (e *Executor) Execute(taskName string, profile string, useCache bool) error {
	if profile != "" {
		e.applyProfile(profile)
	}

	if err := vars.ResolveVars(e.vars); err != nil {
		return err
	}

	task, err := e.graph.GetTask(taskName)
	if err != nil {
		return err
	}

	deps, err := e.graph.GetDependencies(taskName)
	if err != nil {
		return err
	}

	if task.Parallel && len(deps) > 0 {
		if err := e.executeDependenciesParallel(deps, useCache); err != nil {
			return err
		}
	} else {
		for _, dep := range deps {
			depTask, err := e.graph.GetTask(dep)
			if err != nil {
				return err
			}
			if err := e.executeTask(depTask, useCache); err != nil {
				return err
			}
		}
	}

	return e.executeTask(task, useCache)
}

func (e *Executor) executeTask(task *ast.Task, useCache bool) error {
	e.logger.TaskStart(task.Name)
	start := time.Now()

	if e.logStore == nil {
		e.logStore, _ = logs.NewLogStore(logs.GetLogDir())
	}
	if e.logStore != nil {
		e.logStore.StartTask(task.Name)
		e.logStore.Log("info", fmt.Sprintf("Starting task: %s", task.Name))
	}

	taskVars := vars.MergeVars(e.vars, task.Env)

	if task.Profile != "" {
		e.applyProfile(task.Profile)
		taskVars = vars.MergeVars(e.vars, task.Env)
	}

	if len(task.Secrets) > 0 {
		if err := e.loadSecrets(task.Secrets, taskVars); err != nil {
			return err
		}
	}

	if task.If != "" {
		shouldRun, err := e.evaluateCondition(task.If, taskVars)
		if err != nil {
			return fmt.Errorf("condition evaluation failed: %w", err)
		}
		if !shouldRun {
			e.logger.Info(fmt.Sprintf("Skipping task %s (condition not  met)", task.Name))
			if e.collector != nil {
				e.collector.AddSkipped(task.Name)
			}
			return nil
		}
	}

	if len(task.Pre) > 0 {
		if err := e.checkPreconditions(task.Pre); err != nil {
			return err
		}
	}

	if task.Prompt != "" {
		fmt.Printf("%s [y/N]: ", task.Prompt)
		var response string
		_, _ = fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			return fmt.Errorf("task aborted by user")
		}
	}

	cached, inputHash := e.checkEnhancedCache(task, useCache)
	if cached {
		e.logger.TaskCached(task.Name)
		if e.collector != nil {
			e.collector.Add(task.Name, 0, true, true, nil)
		}
		return nil
	}

	if !cached && useCache && len(task.Watch) > 0 {
		hash, err := cache.HashFiles(task.Watch)
		if err == nil {
			if entry, ok := e.cache.Get(task.Name, hash); ok && entry.Success {
				e.logger.TaskCached(task.Name)
				if e.collector != nil {
					e.collector.Add(task.Name, 0, true, true, nil)
				}
				return nil
			}
		}
	}

	success := true
	var execErr error

	if task.Timeout != "" || task.Retries > 0 {
		execErr = e.executeWithTimeout(task, taskVars)
		success = (execErr == nil)
	} else {
		expandedRun := vars.ExpandSlice(task.Run, taskVars)
		for _, cmd := range expandedRun {
			if err := e.runCommand(cmd, taskVars); err != nil {
				success = false
				execErr = err
				break
			}
		}
	}

	duration := time.Since(start)

	if success && useCache {
		if task.Cache && len(task.Inputs) > 0 && inputHash != "" {
			entry := &cache.CacheEntry{
				TaskName:  task.Name,
				InputHash: inputHash,
				Success:   success,
				Duration:  duration,
				Timestamp: time.Now(),
			}
			_ = e.cache.Set(entry)
		} else if len(task.Watch) > 0 {
			hash, _ := cache.HashFiles(task.Watch)
			entry := &cache.CacheEntry{
				TaskName:  task.Name,
				InputHash: hash,
				Success:   success,
				Duration:  duration,
				Timestamp: time.Now(),
			}
			_ = e.cache.Set(entry)
		}
	}

	if e.collector != nil {
		e.collector.Add(task.Name, duration, success, false, execErr)
	}

	if success {
		e.logger.TaskComplete(task.Name, duration)
		if task.Notify.Success != "" {
			e.sendNotification("Flux Task Success", task.Notify.Success)
		}
		if e.logStore != nil {
			e.logStore.Log("info", fmt.Sprintf("Task completed in %v", duration))
			e.logStore.EndTask(task.Name, true)
			_ = e.logStore.Save()
		}
	} else {
		e.logger.TaskFailed(task.Name, execErr)
		if task.Notify.Failure != "" {
			e.sendNotification("Flux Task Failure", task.Notify.Failure)
		}
		if e.logStore != nil {
			e.logStore.Log("error", fmt.Sprintf("Task failed: %v", execErr))
			e.logStore.EndTask(task.Name, false)
			_ = e.logStore.Save()
		}
	}

	return execErr
}

func (e *Executor) sendNotification(title, message string) {
	if e.dryRun {
		e.logger.Info(fmt.Sprintf("[DryRun] Notification: %s - %s", title, message))
		return
	}
	// Simple cross-platform notification using 'msg' on Windows or 'notify-send' on Linux/Mac
	// This is a basic implementation
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("msg", "*", fmt.Sprintf("%s: %s", title, message))
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf("display notification \"%s\" with title \"%s\"", message, title))
	default:
		cmd = exec.Command("notify-send", title, message)
	}
	_ = cmd.Start()
}

func (e *Executor) runCommand(command string, env map[string]string) error {
	if e.dryRun {
		e.logger.Info(fmt.Sprintf("[DryRun] %s", command))
		return nil
	}

	e.logger.Command(command)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	cmd.Env = os.Environ()

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go streamOutput(stdout, e.logger.Stdout)
	go streamOutput(stderr, e.logger.Stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func streamOutput(r io.Reader, writer func(string)) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		writer(scanner.Text())
	}
}

func (e *Executor) applyProfile(profileName string) {
	for _, profile := range e.fluxFile.Profiles {
		if profile.Name == profileName {
			e.vars = vars.MergeVars(e.vars, profile.Env)
			e.logger.Info(fmt.Sprintf("Applied profile: %s", profileName))
			return
		}
	}
	e.logger.Warn(fmt.Sprintf("Profile not found: %s", profileName))
}

func (e *Executor) ListTasks() []string {
	var tasks []string
	for _, task := range e.fluxFile.Tasks {
		tasks = append(tasks, task.Name)
	}
	return tasks
}

func (e *Executor) GetTaskInfo(taskName string) (*ast.Task, error) {
	return e.graph.GetTask(taskName)
}

func (e *Executor) ExecuteAll() error {
	order, err := e.graph.TopologicalSort()
	if err != nil {
		return err
	}

	for _, taskName := range order {
		task, err := e.graph.GetTask(taskName)
		if err != nil {
			return err
		}
		if err := e.executeTask(task, false); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) ExpandMatrixTask(task *ast.Task) []ast.Task {
	if task.Matrix == nil || len(task.Matrix.Dimensions) == 0 {
		return []ast.Task{*task}
	}

	var results []ast.Task
	e.expandMatrix(task, make(map[string]string), &results)
	return results
}

func (e *Executor) expandMatrix(task *ast.Task, current map[string]string, results *[]ast.Task) {
	if len(current) == len(task.Matrix.Dimensions) {
		newTask := *task
		newTask.Matrix = nil

		nameParts := []string{task.Name}
		for k, v := range current {
			nameParts = append(nameParts, fmt.Sprintf("%s=%s", k, v))
		}
		newTask.Name = strings.Join(nameParts, "-")

		newTask.Env = vars.MergeVars(task.Env, current)

		*results = append(*results, newTask)
		return
	}

	for dim, values := range task.Matrix.Dimensions {
		if _, ok := current[dim]; ok {
			continue
		}

		for _, val := range values {
			next := make(map[string]string)
			for k, v := range current {
				next[k] = v
			}
			next[dim] = val
			e.expandMatrix(task, next, results)
		}
		break
	}
}
