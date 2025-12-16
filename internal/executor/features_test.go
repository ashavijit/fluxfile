package executor

import (
	"testing"

	"github.com/ashavijit/fluxfile/internal/ast"
)

func TestResolveAlias(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Tasks = append(fluxFile.Tasks, ast.Task{Name: "build", Alias: "b"})
	fluxFile.Aliases["b"] = "build"

	tests := []struct {
		input    string
		expected string
	}{
		{"b", "build"},
		{"build", "build"},
		{"unknown", "unknown"},
	}

	for _, tt := range tests {
		result := ResolveAlias(fluxFile, tt.input)
		if result != tt.expected {
			t.Errorf("ResolveAlias(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetTaskByNameOrAlias(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Tasks = append(fluxFile.Tasks, ast.Task{Name: "build", Alias: "b"})
	fluxFile.Aliases["b"] = "build"

	task := GetTaskByNameOrAlias(fluxFile, "b")
	if task == nil {
		t.Fatal("Expected task, got nil")
	}
	if task.Name != "build" {
		t.Errorf("Expected task name build, got %s", task.Name)
	}

	taskDirect := GetTaskByNameOrAlias(fluxFile, "build")
	if taskDirect == nil {
		t.Fatal("Expected task, got nil")
	}
	if taskDirect.Name != "build" {
		t.Errorf("Expected task name build, got %s", taskDirect.Name)
	}

	taskUnknown := GetTaskByNameOrAlias(fluxFile, "unknown")
	if taskUnknown != nil {
		t.Error("Expected nil for unknown task")
	}
}

func TestApplyTemplate(t *testing.T) {
	template := ast.NewTemplate("base")
	template.Env = map[string]string{"BASE_VAR": "base_value", "OVERRIDE_ME": "base_v"}
	template.Deps = []string{"dep1"}
	template.Cache = true
	template.Before = []string{"echo base_before"}

	task := ast.NewTask("derived")
	task.Env = map[string]string{"TASK_VAR": "task_value", "OVERRIDE_ME": "task_v"}
	// task.Deps is empty, should inherit
	// task.Cache is false (default), should inherit true
	task.Before = []string{"echo task_before"} // Should NOT inherit if defined

	applyTemplate(&task, &template)

	// Check Env merging
	if task.Env["BASE_VAR"] != "base_value" {
		t.Error("Expected BASE_VAR to be inherited")
	}
	if task.Env["TASK_VAR"] != "task_value" {
		t.Error("Expected TASK_VAR to be preserved")
	}
	if task.Env["OVERRIDE_ME"] != "task_v" {
		t.Errorf("Expected OVERRIDE_ME to be task_v, got %s", task.Env["OVERRIDE_ME"])
	}

	// Check Deps inheritance
	if len(task.Deps) != 1 || task.Deps[0] != "dep1" {
		t.Error("Expected Deps to be inherited")
	}

	// Check Cache inheritance
	if !task.Cache {
		t.Error("Expected Cache to be inherited")
	}

	// Check Before inheritance (Task overrides, so it should NOT inherit, but update logic says "if len(task.Before) == 0")
	if len(task.Before) != 1 || task.Before[0] != "echo task_before" {
		t.Error("Expected task's Before to take precedence")
	}
}

func TestExpandTemplates(t *testing.T) {
	fluxFile := ast.NewFluxFile()

	tmpl := ast.NewTemplate("base")
	tmpl.Cache = true
	fluxFile.Templates = append(fluxFile.Templates, tmpl)

	task := ast.NewTask("build")
	task.Extends = "base"
	fluxFile.Tasks = append(fluxFile.Tasks, task)

	ExpandTemplates(fluxFile)

	if !fluxFile.Tasks[0].Cache {
		t.Error("Expected task to inherit Cache=true from template")
	}
}
