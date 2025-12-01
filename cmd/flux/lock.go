package main

import (
	"fmt"

	"github.com/ashavijit/fluxfile/internal/config"
	"github.com/ashavijit/fluxfile/internal/lock"
)

func handleLockCommands(generateLock, checkLock bool, fluxFilePath string) bool {
	if !generateLock && !checkLock {
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

	if generateLock {
		lockFile, err := lock.Generate(fluxFile)
		if err != nil {
			fmt.Printf("[ERROR] Failed to generate lock: %s\n", err.Error())
			return true
		}

		lockPath := "FluxFile.lock"
		if err := lock.Save(lockFile, lockPath); err != nil {
			fmt.Printf("[ERROR] Failed to save lock: %s\n", err.Error())
			return true
		}

		fmt.Printf("[✓] Lock file generated: %s\n", lockPath)
		fmt.Printf("    Tasks locked: %d\n", len(lockFile.Tasks))
		for taskName, taskLock := range lockFile.Tasks {
			fmt.Printf("    - %s (%d inputs, %d outputs)\n",
				taskName, len(taskLock.Inputs), len(taskLock.Outputs))
		}
		return true
	}

	if checkLock {
		lockPath := "FluxFile.lock"
		lockFile, err := lock.Load(lockPath)
		if err != nil {
			fmt.Printf("[ERROR] Failed to load lock: %s\n", err.Error())
			return true
		}

		changes, err := lock.Verify(lockFile)
		if err != nil {
			fmt.Printf("[ERROR] Failed to verify lock: %s\n", err.Error())
			return true
		}

		if len(changes) == 0 {
			fmt.Println("[✓] Lock file verified - all files match")
			return true
		}

		fmt.Printf("[⚠] Lock file verification failed - %d tasks changed:\n", len(changes))
		for taskName, taskChanges := range changes {
			fmt.Printf("\n  Task: %s\n", taskName)
			for _, change := range taskChanges {
				fmt.Printf("    - %s\n", change)
			}
		}
		return true
	}

	return false
}
