package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/lexer"
	"github.com/ashavijit/fluxfile/internal/parser"
)

type FluxConfig struct {
	DefaultProfile string            `json:"default_profile,omitempty"`
	CacheDir       string            `json:"cache_dir,omitempty"`
	LogDir         string            `json:"log_dir,omitempty"`
	Verbosity      string            `json:"verbosity,omitempty"`
	Parallel       bool              `json:"parallel,omitempty"`
	NoCache        bool              `json:"no_cache,omitempty"`
	WatchDebounce  string            `json:"watch_debounce,omitempty"`
	Env            map[string]string `json:"env,omitempty"`
}

func DefaultConfig() *FluxConfig {
	return &FluxConfig{
		CacheDir:      ".flux/cache",
		LogDir:        ".flux/logs",
		Verbosity:     "normal",
		Parallel:      false,
		NoCache:       false,
		WatchDebounce: "100ms",
		Env:           make(map[string]string),
	}
}

func LoadConfig() (*FluxConfig, error) {
	config := DefaultConfig()

	configPaths := []string{
		".fluxconfig",
		".fluxconfig.json",
		".flux/config.json",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
			}

			if err := json.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
			}

			return config, nil
		}
	}

	return config, nil
}

func SaveConfig(config *FluxConfig, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if path == "" {
		path = ".fluxconfig"
	}

	return os.WriteFile(path, data, 0644)
}

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
