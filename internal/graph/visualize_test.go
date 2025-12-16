package graph

import (
	"strings"
	"testing"

	"github.com/ashavijit/fluxfile/internal/ast"
)

func TestRenderASCII(t *testing.T) {
	tasks := []ast.Task{
		{Name: "deploy", Deps: []string{"build", "test"}},
		{Name: "build", Deps: []string{"lint"}},
		{Name: "test", Deps: []string{}},
		{Name: "lint", Deps: []string{}},
	}

	g, err := BuildGraph(tasks)
	if err != nil {
		t.Fatalf("BuildGraph error: %v", err)
	}

	output := g.RenderASCII("")

	if !strings.Contains(output, "deploy") {
		t.Error("Expected output to contain 'deploy'")
	}
	if !strings.Contains(output, "build") {
		t.Error("Expected output to contain 'build'")
	}
	if !strings.Contains(output, "lint") {
		t.Error("Expected output to contain 'lint'")
	}
}

func TestRenderASCIIFiltered(t *testing.T) {
	tasks := []ast.Task{
		{Name: "deploy", Deps: []string{"build"}},
		{Name: "build", Deps: []string{"lint"}},
		{Name: "lint", Deps: []string{}},
		{Name: "other", Deps: []string{}},
	}

	g, err := BuildGraph(tasks)
	if err != nil {
		t.Fatalf("BuildGraph error: %v", err)
	}

	output := g.RenderASCII("build")

	if !strings.Contains(output, "build") {
		t.Error("Expected output to contain 'build'")
	}
	if !strings.Contains(output, "lint") {
		t.Error("Expected output to contain 'lint'")
	}
}

func TestRenderDOT(t *testing.T) {
	tasks := []ast.Task{
		{Name: "build", Deps: []string{"lint"}},
		{Name: "lint", Deps: []string{}},
	}

	g, err := BuildGraph(tasks)
	if err != nil {
		t.Fatalf("BuildGraph error: %v", err)
	}

	output := g.RenderDOT("")

	if !strings.Contains(output, "digraph FluxFile") {
		t.Error("Expected DOT digraph header")
	}
	if !strings.Contains(output, "\"build\" -> \"lint\"") {
		t.Error("Expected edge from build to lint")
	}
	if !strings.HasSuffix(strings.TrimSpace(output), "}") {
		t.Error("Expected DOT to end with closing brace")
	}
}

func TestRenderMermaid(t *testing.T) {
	tasks := []ast.Task{
		{Name: "build", Deps: []string{"lint"}},
		{Name: "lint", Deps: []string{}},
	}

	g, err := BuildGraph(tasks)
	if err != nil {
		t.Fatalf("BuildGraph error: %v", err)
	}

	output := g.RenderMermaid("")

	if !strings.Contains(output, "```mermaid") {
		t.Error("Expected Mermaid code block start")
	}
	if !strings.Contains(output, "graph TD") {
		t.Error("Expected Mermaid graph directive")
	}
	if !strings.Contains(output, "build --> lint") {
		t.Error("Expected edge from build to lint")
	}
	if !strings.Contains(output, "```\n") {
		t.Error("Expected Mermaid code block end")
	}
}

func TestRenderEmptyGraph(t *testing.T) {
	g := New()

	output := g.RenderASCII("")
	if !strings.Contains(output, "No tasks found") {
		t.Error("Expected 'No tasks found' message for empty graph")
	}
}
