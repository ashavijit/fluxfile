package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/lexer"
	"github.com/ashavijit/fluxfile/internal/parser"
)

func Load(path string) (*ast.FluxFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read FluxFile: %w", err)
	}

	l := lexer.New(string(data))
	p := parser.New(l)

	fluxFile, err := p.Parse()
	if err != nil {
		return nil, fmt.Errorf("failed to parse FluxFile: %w", err)
	}

	baseDir := filepath.Dir(path)
	for _, include := range fluxFile.Includes {
		includePath := filepath.Join(baseDir, include)
		includedFile, err := Load(includePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load included file %s: %w", include, err)
		}

		for k, v := range includedFile.Vars {
			if _, exists := fluxFile.Vars[k]; !exists {
				fluxFile.Vars[k] = v
			}
		}

		fluxFile.Tasks = append(fluxFile.Tasks, includedFile.Tasks...)
		fluxFile.Profiles = append(fluxFile.Profiles, includedFile.Profiles...)
	}

	return fluxFile, nil
}

func FindFluxFile() (string, error) {
	candidates := []string{
		"FluxFile",
		"fluxfile",
		"Fluxfile",
	}

	for _, name := range candidates {
		if _, err := os.Stat(name); err == nil {
			return name, nil
		}
	}

	return "", fmt.Errorf("FluxFile not found in current directory")
}
