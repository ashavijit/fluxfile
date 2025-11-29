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
