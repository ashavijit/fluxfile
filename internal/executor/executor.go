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
	"github.com/ashavijit/fluxfile/internal/vars"
)

type Executor struct {
	fluxFile *ast.FluxFile
	graph    *graph.Graph
	cache    *cache.Cache
	logger   *logger.Logger
	vars     map[string]string
}

type ExecutionResult struct {
	TaskName string
	Success  bool
	Duration time.Duration
	Error    error
}

func New(fluxFile *ast.FluxFile, cacheDir string) (*Executor, error) {
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
	}, nil
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
			return nil
		}
	}

	if len(task.Pre) > 0 {
		if err := e.checkPreconditions(task.Pre); err != nil {
			return err
		}
	}

	cached, inputHash := e.checkEnhancedCache(task, useCache)
	if cached {
		e.logger.TaskCached(task.Name)
		return nil
	}

	if !cached && useCache && len(task.Watch) > 0 {
		hash, err := cache.HashFiles(task.Watch)
		if err == nil {
			if entry, ok := e.cache.Get(task.Name, hash); ok && entry.Success {
				e.logger.TaskCached(task.Name)
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
			e.cache.Set(entry)
		} else if len(task.Watch) > 0 {
			hash, _ := cache.HashFiles(task.Watch)
			entry := &cache.CacheEntry{
				TaskName:  task.Name,
				InputHash: hash,
				Success:   success,
				Duration:  duration,
				Timestamp: time.Now(),
			}
			e.cache.Set(entry)
		}
	}

	if success {
		e.logger.TaskComplete(task.Name, duration)
	} else {
		e.logger.TaskFailed(task.Name, execErr)
	}

	return execErr
}

func (e *Executor) runCommand(command string, env map[string]string) error {
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
