package executor

import (
	"github.com/ashavijit/fluxfile/internal/ast"
)


func ResolveAlias(fluxFile *ast.FluxFile, taskName string) string {
	if actualName, ok := fluxFile.Aliases[taskName]; ok {
		return actualName
	}
	return taskName
}


func GetTaskByNameOrAlias(fluxFile *ast.FluxFile, nameOrAlias string) *ast.Task {
	taskName := ResolveAlias(fluxFile, nameOrAlias)

	for i := range fluxFile.Tasks {
		if fluxFile.Tasks[i].Name == taskName {
			return &fluxFile.Tasks[i]
		}
	}

	return nil
}

func ListAliases(fluxFile *ast.FluxFile) map[string]string {
	return fluxFile.Aliases
}
