package main

import (
	"fmt"

	"github.com/ashavijit/fluxfile/internal/config"
	"github.com/ashavijit/fluxfile/internal/graph"
)

func handleGraphCommands(showGraph, graphDot, graphMermaid bool, filterTask, fluxFilePath string) bool {
	if !showGraph {
		return false
	}

	path := fluxFilePath
	if path == "" {
		var err error
		path, err = config.FindFluxFile()
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
			return true
		}
	}

	fluxFile, err := config.Load(path)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
		return true
	}

	g, err := graph.BuildGraph(fluxFile.Tasks)
	if err != nil {
		fmt.Printf("[ERROR] Failed to build graph: %s\n", err.Error())
		return true
	}

	var output string
	switch {
	case graphDot:
		output = g.RenderDOT(filterTask)
	case graphMermaid:
		output = g.RenderMermaid(filterTask)
	default:
		output = g.RenderASCII(filterTask)
	}

	fmt.Print(output)
	return true
}
