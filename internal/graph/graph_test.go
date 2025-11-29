package graph

import (
	"testing"

	"github.com/ashavijit/fluxfile/internal/ast"
)

func TestTopologicalSort(t *testing.T) {
	tasks := []ast.Task{
		{Name: "a", Deps: []string{}},
		{Name: "b", Deps: []string{"a"}},
		{Name: "c", Deps: []string{"b"}},
	}

	g, err := BuildGraph(tasks)
	if err != nil {
		t.Fatalf("BuildGraph error: %v", err)
	}

	order, err := g.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort error: %v", err)
	}

	if len(order) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(order))
	}

	aIndex := -1
	bIndex := -1
	for i, name := range order {
		if name == "a" {
			aIndex = i
		}
		if name == "b" {
			bIndex = i
		}
	}

	if aIndex > bIndex {
		t.Error("Task 'a' should come before 'b'")
	}
}

func TestCircularDependency(t *testing.T) {
	tasks := []ast.Task{
		{Name: "a", Deps: []string{"b"}},
		{Name: "b", Deps: []string{"a"}},
	}

	_, err := BuildGraph(tasks)
	if err == nil {
		t.Error("Expected circular dependency error")
	}
}

func TestUndefinedDependency(t *testing.T) {
	tasks := []ast.Task{
		{Name: "a", Deps: []string{"nonexistent"}},
	}

	_, err := BuildGraph(tasks)
	if err == nil {
		t.Error("Expected undefined dependency error")
	}
}
