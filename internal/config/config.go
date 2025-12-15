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

// FluxConfig represents project-level configuration
type FluxConfig struct {
	// DefaultProfile is the profile to use when none is specified
	DefaultProfile string `json:"default_profile,omitempty"`

	// CacheDir is the directory for task caching (default: .flux/cache)
	CacheDir string `json:"cache_dir,omitempty"`

	// LogDir is the directory for execution logs (default: .flux/logs)
	LogDir string `json:"log_dir,omitempty"`

	// Verbosity controls log output level: "quiet", "normal", "verbose"
	Verbosity string `json:"verbosity,omitempty"`

	// Parallel controls default parallel execution behavior
	Parallel bool `json:"parallel,omitempty"`

	// NoCache disables caching by default
	NoCache bool `json:"no_cache,omitempty"`

	// WatchDebounce is the debounce duration for file watching (e.g., "100ms")
	WatchDebounce string `json:"watch_debounce,omitempty"`

	// Environment variables to set for all tasks
	Env map[string]string `json:"env,omitempty"`
}

// DefaultConfig returns a FluxConfig with sensible defaults
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

// LoadConfig loads configuration from .fluxconfig file
func LoadConfig() (*FluxConfig, error) {
	config := DefaultConfig()

	// Look for config files in order of priority
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

	// No config file found, return defaults
	return config, nil
}

// SaveConfig saves configuration to .fluxconfig file
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
