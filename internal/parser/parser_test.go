package parser

import (
	"testing"

	"github.com/ashavijit/fluxfile/internal/lexer"
)

func TestParseVarDecl(t *testing.T) {
	input := `var PROJECT = flux`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if fluxFile.Vars["PROJECT"] != "flux" {
		t.Errorf("Expected PROJECT=flux, got %s", fluxFile.Vars["PROJECT"])
	}
}

func TestParseTask(t *testing.T) {
	input := `task build:
    run:
        go build
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(fluxFile.Tasks))
	}

	if fluxFile.Tasks[0].Name != "build" {
		t.Errorf("Expected task name build, got %s", fluxFile.Tasks[0].Name)
	}
}

func TestParseDeps(t *testing.T) {
	input := `task deploy:
    deps: build, test
    run:
        ./deploy.sh
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Tasks[0].Deps) != 2 {
		t.Errorf("Expected 2 deps, got %d", len(fluxFile.Tasks[0].Deps))
	}
}

func TestParseProfile(t *testing.T) {
	input := `profile dev:
    env:
        MODE = development
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Profiles) != 1 {
		t.Fatalf("Expected 1 profile, got %d", len(fluxFile.Profiles))
	}

	if fluxFile.Profiles[0].Name != "dev" {
		t.Errorf("Expected profile name dev, got %s", fluxFile.Profiles[0].Name)
	}
}

func TestParseRetries(t *testing.T) {
	input := `task flaky:
    retries: 3
    run:
        ./flaky.sh
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if fluxFile.Tasks[0].Retries != 3 {
		t.Errorf("Expected retries 3, got %d", fluxFile.Tasks[0].Retries)
	}
}

func TestParseAlias(t *testing.T) {
	input := `task build:
    alias: b
    run:
        go build
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if fluxFile.Tasks[0].Alias != "b" {
		t.Errorf("Expected alias b, got %s", fluxFile.Tasks[0].Alias)
	}

	if fluxFile.Aliases["b"] != "build" {
		t.Errorf("Expected alias mapping b->build, got %s", fluxFile.Aliases["b"])
	}
}

func TestParseExtends(t *testing.T) {
	input := `task build:
    extends: go-base
    run:
        go build
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if fluxFile.Tasks[0].Extends != "go-base" {
		t.Errorf("Expected extends go-base, got %s", fluxFile.Tasks[0].Extends)
	}
}

func TestParseBeforeAfterHooks(t *testing.T) {
	input := `task deploy:
    before:
        echo starting
    run:
        ./deploy.sh
    after:
        echo done
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Tasks[0].Before) != 1 {
		t.Errorf("Expected 1 before hook, got %d", len(fluxFile.Tasks[0].Before))
	}

	if len(fluxFile.Tasks[0].After) != 1 {
		t.Errorf("Expected 1 after hook, got %d", len(fluxFile.Tasks[0].After))
	}
}

func TestParseTemplate(t *testing.T) {
	input := `template go-base:
    cache: true
    inputs:
        **/*.go
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Templates) != 1 {
		t.Fatalf("Expected 1 template, got %d", len(fluxFile.Templates))
	}

	if fluxFile.Templates[0].Name != "go-base" {
		t.Errorf("Expected template name go-base, got %s", fluxFile.Templates[0].Name)
	}

	if !fluxFile.Templates[0].Cache {
		t.Error("Expected template cache to be true")
	}
}

func TestParseTaskGroup(t *testing.T) {
	input := `group frontend:
    tasks: install, build, test
`

	l := lexer.New(input)
	p := New(l)
	fluxFile, err := p.Parse()

	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(fluxFile.Groups) != 1 {
		t.Fatalf("Expected 1 group, got %d", len(fluxFile.Groups))
	}

	if fluxFile.Groups[0].Name != "frontend" {
		t.Errorf("Expected group name frontend, got %s", fluxFile.Groups[0].Name)
	}

	if len(fluxFile.Groups[0].Tasks) != 3 {
		t.Errorf("Expected 3 tasks in group, got %d", len(fluxFile.Groups[0].Tasks))
	}
}
