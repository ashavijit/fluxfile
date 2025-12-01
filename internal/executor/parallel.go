package executor

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

func (e *Executor) executeDependenciesParallel(deps []string, useCache bool) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(deps))

	for _, dep := range deps {
		wg.Add(1)
		go func(depName string) {
			defer wg.Done()
			depTask, err := e.graph.GetTask(depName)
			if err != nil {
				errors <- err
				return
			}
			if err := e.executeTask(depTask, useCache); err != nil {
				errors <- fmt.Errorf("dependency %s failed: %w", depName, err)
			}
		}(dep)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) loadSecrets(secrets []string, vars map[string]string) error {
	for _, secretKey := range secrets {
		value := os.Getenv(secretKey)
		if value == "" {
			envFile := ".env"
			if envContent, err := os.ReadFile(envFile); err == nil {
				lines := strings.Split(string(envContent), "\n")
				for _, line := range lines {
					line = strings.TrimSpace(line)
					if line == "" || strings.HasPrefix(line, "#") {
						continue
					}
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 && parts[0] == secretKey {
						value = strings.Trim(parts[1], "\"'")
						break
					}
				}
			}
		}

		if value != "" {
			vars[secretKey] = value
		} else {
			return fmt.Errorf("secret %s not found in environment or .env file", secretKey)
		}
	}
	return nil
}
