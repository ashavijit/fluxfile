package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("Expected config, got nil")
	}

	if config.CacheDir != ".flux/cache" {
		t.Errorf("Expected CacheDir .flux/cache, got %s", config.CacheDir)
	}

	if config.LogDir != ".flux/logs" {
		t.Errorf("Expected LogDir .flux/logs, got %s", config.LogDir)
	}

	if config.Verbosity != "normal" {
		t.Errorf("Expected Verbosity normal, got %s", config.Verbosity)
	}

	if config.WatchDebounce != "100ms" {
		t.Errorf("Expected WatchDebounce 100ms, got %s", config.WatchDebounce)
	}

	if config.Parallel {
		t.Error("Expected Parallel to be false by default")
	}

	if config.NoCache {
		t.Error("Expected NoCache to be false by default")
	}

	if config.Env == nil {
		t.Error("Expected Env to be initialized")
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	// Change to temp dir where no config exists
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig should not error when no config exists: %v", err)
	}

	// Should return defaults when no config file exists
	if config.CacheDir != ".flux/cache" {
		t.Errorf("Expected default CacheDir, got %s", config.CacheDir)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	// Create a config file
	configContent := `{
  "default_profile": "prod",
  "cache_dir": "custom/cache",
  "verbosity": "verbose",
  "parallel": true,
  "env": {
    "APP_ENV": "production"
  }
}`
	os.WriteFile(".fluxconfig", []byte(configContent), 0644)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.DefaultProfile != "prod" {
		t.Errorf("Expected DefaultProfile prod, got %s", config.DefaultProfile)
	}

	if config.CacheDir != "custom/cache" {
		t.Errorf("Expected CacheDir custom/cache, got %s", config.CacheDir)
	}

	if config.Verbosity != "verbose" {
		t.Errorf("Expected Verbosity verbose, got %s", config.Verbosity)
	}

	if !config.Parallel {
		t.Error("Expected Parallel to be true")
	}

	if config.Env["APP_ENV"] != "production" {
		t.Errorf("Expected APP_ENV=production, got %s", config.Env["APP_ENV"])
	}
}

func TestLoadConfigJsonExtension(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	// Create a config file with .json extension
	configContent := `{"verbosity": "quiet"}`
	os.WriteFile(".fluxconfig.json", []byte(configContent), 0644)

	config, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if config.Verbosity != "quiet" {
		t.Errorf("Expected Verbosity quiet, got %s", config.Verbosity)
	}
}

func TestLoadConfigInvalidJson(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	// Create an invalid config file
	os.WriteFile(".fluxconfig", []byte("invalid json {"), 0644)

	_, err := LoadConfig()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestSaveConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".fluxconfig")

	config := &FluxConfig{
		DefaultProfile: "dev",
		CacheDir:       "my-cache",
		Verbosity:      "verbose",
		Parallel:       true,
		Env: map[string]string{
			"KEY": "value",
		},
	}

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load and verify
	data, _ := os.ReadFile(configPath)
	if len(data) == 0 {
		t.Error("Config file is empty")
	}
}

func TestSaveConfigDefaultPath(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	config := DefaultConfig()
	err := SaveConfig(config, "")
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// Should save to .fluxconfig by default
	if _, err := os.Stat(".fluxconfig"); os.IsNotExist(err) {
		t.Fatal("Config file was not created at default path")
	}
}

func TestFindFluxFile(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	// No FluxFile exists
	_, err := FindFluxFile()
	if err == nil {
		t.Error("Expected error when FluxFile doesn't exist")
	}

	// Create FluxFile
	os.WriteFile("FluxFile", []byte("task build:\n    run:\n        echo hi"), 0644)

	path, err := FindFluxFile()
	if err != nil {
		t.Fatalf("FindFluxFile failed: %v", err)
	}

	if path != "FluxFile" {
		t.Errorf("Expected FluxFile, got %s", path)
	}
}

func TestFindFluxFileCaseVariants(t *testing.T) {
	dir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(dir)

	os.WriteFile("FluxFile", []byte("task test:\n    run:\n        echo test"), 0644)

	path, err := FindFluxFile()
	if err != nil {
		t.Fatalf("FindFluxFile failed: %v", err)
	}

	if path != "FluxFile" {
		t.Errorf("Expected FluxFile, got %s", path)
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	fluxFilePath := filepath.Join(dir, "FluxFile")

	content := `var PROJECT = test

task build:
    run:
        go build

task test:
    deps: build
    run:
        go test
`
	os.WriteFile(fluxFilePath, []byte(content), 0644)

	fluxFile, err := Load(fluxFilePath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if fluxFile.Vars["PROJECT"] != "test" {
		t.Errorf("Expected PROJECT=test, got %s", fluxFile.Vars["PROJECT"])
	}

	if len(fluxFile.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(fluxFile.Tasks))
	}
}

func TestLoadNotFound(t *testing.T) {
	_, err := Load("/nonexistent/FluxFile")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
