package executor

import (
	"github.com/ashavijit/fluxfile/internal/ast"
)


func ExpandTemplates(fluxFile *ast.FluxFile) {
	templateMap := make(map[string]*ast.Template)
	for i := range fluxFile.Templates {
		templateMap[fluxFile.Templates[i].Name] = &fluxFile.Templates[i]
	}

	for i := range fluxFile.Tasks {
		task := &fluxFile.Tasks[i]
		if task.Extends != "" {
			if template, ok := templateMap[task.Extends]; ok {
				applyTemplate(task, template)
			}
		}
	}
}

// applyTemplate merges template properties into a task.
// Task-specific values take precedence over template values.
func applyTemplate(task *ast.Task, template *ast.Template) {
	if task.Desc == "" {
		task.Desc = template.Desc
	}

	if len(task.Deps) == 0 {
		task.Deps = template.Deps
	}

	if len(task.Env) == 0 {
		task.Env = template.Env
	} else if len(template.Env) > 0 {
		merged := make(map[string]string)
		for k, v := range template.Env {
			merged[k] = v
		}
		for k, v := range task.Env {
			merged[k] = v
		}
		task.Env = merged
	}

	if !task.Cache && template.Cache {
		task.Cache = template.Cache
	}

	if len(task.Inputs) == 0 {
		task.Inputs = template.Inputs
	}

	if len(task.Outputs) == 0 {
		task.Outputs = template.Outputs
	}

	if !task.Parallel && template.Parallel {
		task.Parallel = template.Parallel
	}

	if !task.Docker && template.Docker {
		task.Docker = template.Docker
	}

	if task.Remote == "" {
		task.Remote = template.Remote
	}

	if len(task.Secrets) == 0 {
		task.Secrets = template.Secrets
	}

	if len(task.Pre) == 0 {
		task.Pre = template.Pre
	}

	if task.Retries == 0 {
		task.Retries = template.Retries
	}

	if task.RetryDelay == "" {
		task.RetryDelay = template.RetryDelay
	}

	if task.Timeout == "" {
		task.Timeout = template.Timeout
	}

	if len(task.Before) == 0 {
		task.Before = template.Before
	}

	if len(task.After) == 0 {
		task.After = template.After
	}
}
