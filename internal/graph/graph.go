package graph

import (
	"fmt"

	"github.com/ashavijit/fluxfile/internal/ast"
)

type Graph struct {
	tasks map[string]*ast.Task
	edges map[string][]string
}

func New() *Graph {
	return &Graph{
		tasks: make(map[string]*ast.Task),
		edges: make(map[string][]string),
	}
}

func (g *Graph) AddTask(task *ast.Task) {
	g.tasks[task.Name] = task
	g.edges[task.Name] = task.Deps
}

func (g *Graph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	var result []string

	var visit func(string) error
	visit = func(node string) error {
		if visited[node] {
			return nil
		}
		if recStack[node] {
			return fmt.Errorf("circular dependency detected: %s", node)
		}

		recStack[node] = true
		for _, dep := range g.edges[node] {
			if _, exists := g.tasks[dep]; !exists {
				return fmt.Errorf("task %s depends on undefined task %s", node, dep)
			}
			if err := visit(dep); err != nil {
				return err
			}
		}
		delete(recStack, node)
		visited[node] = true
		result = append(result, node)
		return nil
	}

	for name := range g.tasks {
		if !visited[name] {
			if err := visit(name); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

func (g *Graph) GetTask(name string) (*ast.Task, error) {
	task, ok := g.tasks[name]
	if !ok {
		return nil, fmt.Errorf("task %s not found", name)
	}
	return task, nil
}

func (g *Graph) GetDependencies(taskName string) ([]string, error) {
	task, err := g.GetTask(taskName)
	if err != nil {
		return nil, err
	}

	visited := make(map[string]bool)
	var result []string

	var collect func(string) error
	collect = func(name string) error {
		if visited[name] {
			return nil
		}
		visited[name] = true

		t, err := g.GetTask(name)
		if err != nil {
			return err
		}

		for _, dep := range t.Deps {
			if err := collect(dep); err != nil {
				return err
			}
		}

		if name != taskName {
			result = append(result, name)
		}
		return nil
	}

	for _, dep := range task.Deps {
		if err := collect(dep); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func BuildGraph(tasks []ast.Task) (*Graph, error) {
	g := New()
	for i := range tasks {
		g.AddTask(&tasks[i])
	}

	_, err := g.TopologicalSort()
	if err != nil {
		return nil, err
	}

	return g, nil
}
