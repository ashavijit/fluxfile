package vars

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var shellExprPattern = regexp.MustCompile(`\$\(shell\s+"([^"]+)"\)`)
var varPattern = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_-]*)\}`)

func Expand(value string, vars map[string]string) string {
	result := value

	result = shellExprPattern.ReplaceAllStringFunc(result, func(match string) string {
		matches := shellExprPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		command := matches[1]
		return executeShellCommand(command)
	})

	result = varPattern.ReplaceAllStringFunc(result, func(match string) string {
		matches := varPattern.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		varName := matches[1]
		if val, ok := vars[varName]; ok {
			return Expand(val, vars)
		}
		if val, ok := os.LookupEnv(varName); ok {
			return val
		}
		return match
	})

	return result
}

func executeShellCommand(command string) string {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell.exe", "-Command", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func ExpandMap(m map[string]string, vars map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = Expand(v, vars)
	}
	return result
}

func ExpandSlice(s []string, vars map[string]string) []string {
	result := make([]string, len(s))
	for i, v := range s {
		result[i] = Expand(v, vars)
	}
	return result
}

func MergeVars(base, overlay map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		result[k] = v
	}
	return result
}

func ResolveVars(vars map[string]string) error {
	maxIterations := 100
	for i := 0; i < maxIterations; i++ {
		changed := false
		for k, v := range vars {
			expanded := Expand(v, vars)
			if expanded != v {
				vars[k] = expanded
				changed = true
			}
		}
		if !changed {
			return nil
		}
	}
	return fmt.Errorf("circular variable dependency detected")
}
