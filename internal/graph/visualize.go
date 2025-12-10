package graph

import (
	"fmt"
	"sort"
	"strings"
)

// RenderASCII renders the dependency graph as an ASCII tree
func (g *Graph) RenderASCII(filterTask string) string {
	var sb strings.Builder

	// Get all root tasks (tasks with no reverse dependencies or the filtered task)
	roots := g.getRoots(filterTask)
	if len(roots) == 0 {
		return "No tasks found\n"
	}

	sort.Strings(roots)
	for i, root := range roots {
		isLast := i == len(roots)-1
		g.renderASCIINode(&sb, root, "", isLast, make(map[string]bool))
	}

	return sb.String()
}

func (g *Graph) renderASCIINode(sb *strings.Builder, name string, prefix string, isLast bool, visited map[string]bool) {
	// Choose the connector
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	// Print current node
	if prefix == "" {
		sb.WriteString(name)
	} else {
		sb.WriteString(prefix + connector + name)
	}

	// Add task description if available
	if task, ok := g.tasks[name]; ok && task.Desc != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", task.Desc))
	}
	sb.WriteString("\n")

	// Prevent infinite loops
	if visited[name] {
		return
	}
	visited[name] = true

	// Get dependencies
	deps := g.edges[name]
	if len(deps) == 0 {
		return
	}

	// Calculate new prefix for children
	newPrefix := prefix
	if prefix == "" {
		newPrefix = ""
	} else if isLast {
		newPrefix = prefix + "    "
	} else {
		newPrefix = prefix + "│   "
	}

	// Sort dependencies for consistent output
	sortedDeps := make([]string, len(deps))
	copy(sortedDeps, deps)
	sort.Strings(sortedDeps)

	for i, dep := range sortedDeps {
		childIsLast := i == len(sortedDeps)-1
		g.renderASCIINode(sb, dep, newPrefix, childIsLast, visited)
	}
}

// RenderDOT renders the dependency graph in Graphviz DOT format
func (g *Graph) RenderDOT(filterTask string) string {
	var sb strings.Builder

	sb.WriteString("digraph FluxFile {\n")
	sb.WriteString("  rankdir=TB;\n")
	sb.WriteString("  node [shape=box, style=rounded, fontname=\"Arial\"];\n")
	sb.WriteString("  edge [arrowhead=vee];\n")
	sb.WriteString("\n")

	// Collect edges
	edges := g.collectEdges(filterTask)

	// Sort edges for consistent output
	sort.Slice(edges, func(i, j int) bool {
		if edges[i][0] == edges[j][0] {
			return edges[i][1] < edges[j][1]
		}
		return edges[i][0] < edges[j][0]
	})

	// Write edges
	for _, edge := range edges {
		sb.WriteString(fmt.Sprintf("  \"%s\" -> \"%s\";\n", edge[0], edge[1]))
	}

	// Add isolated nodes (no dependencies)
	nodes := g.collectNodes(filterTask)
	for _, node := range nodes {
		hasEdge := false
		for _, edge := range edges {
			if edge[0] == node || edge[1] == node {
				hasEdge = true
				break
			}
		}
		if !hasEdge {
			sb.WriteString(fmt.Sprintf("  \"%s\";\n", node))
		}
	}

	sb.WriteString("}\n")
	return sb.String()
}

// RenderMermaid renders the dependency graph in Mermaid format
func (g *Graph) RenderMermaid(filterTask string) string {
	var sb strings.Builder

	sb.WriteString("```mermaid\n")
	sb.WriteString("graph TD\n")

	// Collect edges
	edges := g.collectEdges(filterTask)

	// Sort edges for consistent output
	sort.Slice(edges, func(i, j int) bool {
		if edges[i][0] == edges[j][0] {
			return edges[i][1] < edges[j][1]
		}
		return edges[i][0] < edges[j][0]
	})

	// Write edges
	for _, edge := range edges {
		sb.WriteString(fmt.Sprintf("  %s --> %s\n", edge[0], edge[1]))
	}

	// Add isolated nodes
	nodes := g.collectNodes(filterTask)
	for _, node := range nodes {
		hasEdge := false
		for _, edge := range edges {
			if edge[0] == node || edge[1] == node {
				hasEdge = true
				break
			}
		}
		if !hasEdge {
			sb.WriteString(fmt.Sprintf("  %s\n", node))
		}
	}

	sb.WriteString("```\n")
	return sb.String()
}

// Helper: get root tasks (no reverse dependencies) or filter to specific task
func (g *Graph) getRoots(filterTask string) []string {
	if filterTask != "" {
		if _, exists := g.tasks[filterTask]; exists {
			return []string{filterTask}
		}
		return nil
	}

	// Find tasks that are not dependencies of any other task
	isDep := make(map[string]bool)
	for _, deps := range g.edges {
		for _, dep := range deps {
			isDep[dep] = true
		}
	}

	var roots []string
	for name := range g.tasks {
		if !isDep[name] {
			roots = append(roots, name)
		}
	}

	// If everything is a dependency of something, return all tasks
	if len(roots) == 0 {
		for name := range g.tasks {
			roots = append(roots, name)
		}
	}

	return roots
}

// Helper: collect all edges for rendering
func (g *Graph) collectEdges(filterTask string) [][2]string {
	var edges [][2]string
	visited := make(map[string]bool)

	var collect func(string)
	collect = func(name string) {
		if visited[name] {
			return
		}
		visited[name] = true

		for _, dep := range g.edges[name] {
			edges = append(edges, [2]string{name, dep})
			collect(dep)
		}
	}

	if filterTask != "" {
		collect(filterTask)
	} else {
		for name := range g.tasks {
			collect(name)
		}
	}

	return edges
}

// Helper: collect all nodes for rendering
func (g *Graph) collectNodes(filterTask string) []string {
	if filterTask != "" {
		visited := make(map[string]bool)
		var collect func(string)
		collect = func(name string) {
			if visited[name] {
				return
			}
			visited[name] = true
			for _, dep := range g.edges[name] {
				collect(dep)
			}
		}
		collect(filterTask)

		var nodes []string
		for name := range visited {
			nodes = append(nodes, name)
		}
		sort.Strings(nodes)
		return nodes
	}

	var nodes []string
	for name := range g.tasks {
		nodes = append(nodes, name)
	}
	sort.Strings(nodes)
	return nodes
}
