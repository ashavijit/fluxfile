package init

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ProjectName string
	Template    string
	Directory   string
}

type ProjectType string

const (
	TypeGo      ProjectType = "go"
	TypeNode    ProjectType = "node"
	TypePython  ProjectType = "python"
	TypeRust    ProjectType = "rust"
	TypeGeneric ProjectType = "generic"
)

func Run(cfg Config) error {
	if cfg.Directory == "" {
		cfg.Directory = "."
	}

	if cfg.ProjectName == "" {
		cfg.ProjectName = detectProjectName(cfg.Directory)
	}

	template := cfg.Template
	if template == "" {
		template = string(Detect(cfg.Directory))
	}

	content := GetTemplate(template, cfg.ProjectName)
	fluxPath := filepath.Join(cfg.Directory, "FluxFile")

	if _, err := os.Stat(fluxPath); err == nil {
		fmt.Print("FluxFile already exists. Overwrite? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			return fmt.Errorf("aborted: FluxFile already exists")
		}
	}

	if err := os.WriteFile(fluxPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write FluxFile: %w", err)
	}

	fluxDir := filepath.Join(cfg.Directory, ".flux")
	if err := os.MkdirAll(fluxDir, 0755); err != nil {
		return fmt.Errorf("failed to create .flux directory: %w", err)
	}

	fmt.Printf("Created FluxFile with %s template\n", template)
	fmt.Printf("Project: %s\n", cfg.ProjectName)
	fmt.Println("\nRun 'flux -l' to list available tasks")

	return nil
}

func Detect(dir string) ProjectType {
	checks := []struct {
		file     string
		projType ProjectType
	}{
		{"go.mod", TypeGo},
		{"package.json", TypeNode},
		{"pyproject.toml", TypePython},
		{"requirements.txt", TypePython},
		{"Cargo.toml", TypeRust},
	}

	for _, check := range checks {
		path := filepath.Join(dir, check.file)
		if _, err := os.Stat(path); err == nil {
			return check.projType
		}
	}

	return TypeGeneric
}

func GetTemplates() []string {
	return []string{"go", "node", "python", "rust", "generic"}
}

func detectProjectName(dir string) string {
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "project"
	}
	return filepath.Base(absPath)
}
