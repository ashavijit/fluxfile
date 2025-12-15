package executor

import (
	"testing"
	"time"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/graph"
)

func TestNew(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Tasks = append(fluxFile.Tasks, ast.NewTask("build"))

	exec, err := New(fluxFile, t.TempDir(), false)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	if exec == nil {
		t.Fatal("Expected executor, got nil")
	}

	if exec.fluxFile != fluxFile {
		t.Error("Expected fluxFile to be set")
	}
}

func TestListTasks(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Tasks = []ast.Task{
		ast.NewTask("build"),
		ast.NewTask("test"),
		ast.NewTask("deploy"),
	}

	exec, err := New(fluxFile, t.TempDir(), false)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	tasks := exec.ListTasks()
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}

	expected := map[string]bool{"build": true, "test": true, "deploy": true}
	for _, task := range tasks {
		if !expected[task] {
			t.Errorf("Unexpected task: %s", task)
		}
	}
}

func TestEvaluateCondition(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	exec, _ := New(fluxFile, t.TempDir(), false)

	tests := []struct {
		name      string
		condition string
		vars      map[string]string
		expected  bool
		wantErr   bool
	}{
		{
			name:      "empty condition",
			condition: "",
			vars:      map[string]string{},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "equal condition true",
			condition: "MODE == prod",
			vars:      map[string]string{"MODE": "prod"},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "equal condition false",
			condition: "MODE == prod",
			vars:      map[string]string{"MODE": "dev"},
			expected:  false,
			wantErr:   false,
		},
		{
			name:      "not equal condition",
			condition: "MODE != prod",
			vars:      map[string]string{"MODE": "dev"},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "greater than numeric",
			condition: "COUNT > 5",
			vars:      map[string]string{"COUNT": "10"},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "less than numeric",
			condition: "COUNT < 5",
			vars:      map[string]string{"COUNT": "3"},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "greater equal numeric",
			condition: "COUNT >= 5",
			vars:      map[string]string{"COUNT": "5"},
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "no operator",
			condition: "INVALID",
			vars:      map[string]string{},
			expected:  false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := exec.evaluateCondition(tt.condition, tt.vars)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("evaluateCondition() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCheckPreconditions(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	exec, _ := New(fluxFile, t.TempDir(), false)

	tests := []struct {
		name       string
		conditions []ast.Precondition
		wantErr    bool
	}{
		{
			name:       "empty preconditions",
			conditions: []ast.Precondition{},
			wantErr:    false,
		},
		{
			name: "file exists - go.mod",
			conditions: []ast.Precondition{
				{Type: "file", Value: "go.mod"},
			},
			wantErr: true, // Will fail because we're in temp dir
		},
		{
			name: "command exists - go",
			conditions: []ast.Precondition{
				{Type: "command", Value: "go"},
			},
			wantErr: false, // go should be installed
		},
		{
			name: "env var set - PATH",
			conditions: []ast.Precondition{
				{Type: "env", Value: "PATH"},
			},
			wantErr: false, // PATH should always be set
		},
		{
			name: "env var not set",
			conditions: []ast.Precondition{
				{Type: "env", Value: "FLUX_TEST_NONEXISTENT_VAR"},
			},
			wantErr: true,
		},
		{
			name: "unknown precondition type",
			conditions: []ast.Precondition{
				{Type: "unknown", Value: "value"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exec.checkPreconditions(tt.conditions)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPreconditions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseRetryDelay(t *testing.T) {
	tests := []struct {
		input    string
		expected time.Duration
	}{
		{"", 1 * time.Second},
		{"2s", 2 * time.Second},
		{"500ms", 500 * time.Millisecond},
		{"1m", 1 * time.Minute},
		{"invalid", 1 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseRetryDelay(tt.input)
			if result != tt.expected {
				t.Errorf("parseRetryDelay(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpandString(t *testing.T) {
	vars := map[string]string{
		"NAME":    "flux",
		"VERSION": "1.0.0",
	}

	tests := []struct {
		input    string
		expected string
	}{
		{"hello ${NAME}", "hello flux"},
		{"version: $VERSION", "version: 1.0.0"},
		{"${NAME}-${VERSION}", "flux-1.0.0"},
		{"no vars here", "no vars here"},
		{"$UNKNOWN", "$UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := expandString(tt.input, vars)
			if result != tt.expected {
				t.Errorf("expandString(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExpandSlice(t *testing.T) {
	vars := map[string]string{
		"DIR": "build",
	}

	input := []string{"go build -o ${DIR}/app", "echo $DIR"}
	result := expandSlice(input, vars)

	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result))
	}

	if result[0] != "go build -o build/app" {
		t.Errorf("Expected expanded command, got %q", result[0])
	}

	if result[1] != "echo build" {
		t.Errorf("Expected expanded echo, got %q", result[1])
	}
}

func TestExecLookPath(t *testing.T) {
	// Test that go command can be found
	path, err := execLookPath("go")
	if err != nil {
		t.Errorf("execLookPath(go) failed: %v", err)
	}
	if path == "" {
		t.Error("Expected non-empty path for go")
	}

	// Test that non-existent command returns error
	_, err = execLookPath("flux_definitely_not_a_real_command_12345")
	if err == nil {
		t.Error("Expected error for non-existent command")
	}
}

func TestExecutorWithGraph(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	buildTask := ast.NewTask("build")
	buildTask.Run = []string{"echo building"}

	testTask := ast.NewTask("test")
	testTask.Deps = []string{"build"}
	testTask.Run = []string{"echo testing"}

	fluxFile.Tasks = []ast.Task{buildTask, testTask}

	exec, err := New(fluxFile, t.TempDir(), false)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	// Verify graph was built
	if exec.graph == nil {
		t.Fatal("Expected graph to be built")
	}

	// Test GetTaskInfo
	taskInfo, err := exec.GetTaskInfo("build")
	if err != nil {
		t.Errorf("GetTaskInfo(build) failed: %v", err)
	}
	if taskInfo == nil {
		t.Error("Expected task info, got nil")
	}
}

func TestDryRun(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	task := ast.NewTask("build")
	task.Run = []string{"echo hello"}
	fluxFile.Tasks = []ast.Task{task}

	exec, err := New(fluxFile, t.TempDir(), true) // dryRun = true
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	if !exec.dryRun {
		t.Error("Expected dryRun to be true")
	}
}

func TestMatrixExpansion(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	exec, _ := New(fluxFile, t.TempDir(), false)

	task := ast.NewTask("build")
	task.Matrix = &ast.Matrix{
		Dimensions: map[string][]string{
			"os":   {"linux", "darwin"},
			"arch": {"amd64", "arm64"},
		},
	}
	task.Run = []string{"echo ${os}-${arch}"}

	expanded := exec.ExpandMatrixTask(&task)

	// 2 OS * 2 arch = 4 combinations
	if len(expanded) != 4 {
		t.Errorf("Expected 4 expanded tasks, got %d", len(expanded))
	}
}

func TestApplyProfile(t *testing.T) {
	fluxFile := ast.NewFluxFile()
	fluxFile.Profiles = []ast.Profile{
		{
			Name: "dev",
			Env: map[string]string{
				"MODE":  "development",
				"DEBUG": "true",
			},
		},
	}

	exec, err := New(fluxFile, t.TempDir(), false)
	if err != nil {
		t.Fatalf("Failed to create executor: %v", err)
	}

	exec.applyProfile("dev")

	if exec.vars["MODE"] != "development" {
		t.Errorf("Expected MODE=development, got %s", exec.vars["MODE"])
	}

	if exec.vars["DEBUG"] != "true" {
		t.Errorf("Expected DEBUG=true, got %s", exec.vars["DEBUG"])
	}
}

func TestGraphIntegration(t *testing.T) {
	tasks := []ast.Task{
		{Name: "clean", Run: []string{"echo clean"}},
		{Name: "build", Deps: []string{"clean"}, Run: []string{"echo build"}},
		{Name: "test", Deps: []string{"build"}, Run: []string{"echo test"}},
		{Name: "deploy", Deps: []string{"test"}, Run: []string{"echo deploy"}},
	}

	g, err := graph.BuildGraph(tasks)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	if g == nil {
		t.Fatal("Expected graph, got nil")
	}

	task, err := g.GetTask("build")
	if err != nil {
		t.Errorf("GetTask(build) failed: %v", err)
	}
	if task == nil {
		t.Error("Expected task, got nil")
	}

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort failed: %v", err)
	}

	if len(order) != 4 {
		t.Errorf("Expected 4 tasks in order, got %d", len(order))
	}
}
