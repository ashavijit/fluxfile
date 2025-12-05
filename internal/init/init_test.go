package init

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected ProjectType
	}{
		{"go project", []string{"go.mod"}, TypeGo},
		{"node project", []string{"package.json"}, TypeNode},
		{"python project", []string{"pyproject.toml"}, TypePython},
		{"rust project", []string{"Cargo.toml"}, TypeRust},
		{"generic project", []string{}, TypeGeneric},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tt.files {
				path := filepath.Join(dir, f)
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatal(err)
				}
			}

			got := Detect(dir)
			if got != tt.expected {
				t.Errorf("Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGetTemplate(t *testing.T) {
	templates := GetTemplates()
	if len(templates) != 5 {
		t.Errorf("Expected 5 templates, got %d", len(templates))
	}

	for _, tmpl := range templates {
		content := GetTemplate(tmpl, "testproj")
		if content == "" {
			t.Errorf("Template %s returned empty content", tmpl)
		}
		if !strings.Contains(content, "testproj") {
			t.Errorf("Template %s does not contain project name", tmpl)
		}
		if !strings.Contains(content, "task ") {
			t.Errorf("Template %s does not contain task definitions", tmpl)
		}
	}
}

func TestRun(t *testing.T) {
	dir := t.TempDir()

	cfg := Config{
		ProjectName: "myproject",
		Template:    "generic",
		Directory:   dir,
	}

	if err := Run(cfg); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	fluxPath := filepath.Join(dir, "FluxFile")
	if _, err := os.Stat(fluxPath); os.IsNotExist(err) {
		t.Error("FluxFile was not created")
	}

	content, err := os.ReadFile(fluxPath)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "myproject") {
		t.Error("FluxFile does not contain project name")
	}

	fluxDir := filepath.Join(dir, ".flux")
	if _, err := os.Stat(fluxDir); os.IsNotExist(err) {
		t.Error(".flux directory was not created")
	}
}
