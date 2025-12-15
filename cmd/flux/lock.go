package main

import (
	"encoding/json"
	"fmt"

	"github.com/ashavijit/fluxfile/internal/config"
	"github.com/ashavijit/fluxfile/internal/lock"
)

func handleLockCommands(generateLock bool, checkLock bool, lockUpdate bool, lockDiff bool, lockClean bool, updateTask string, fluxFilePath string, jsonOutput bool) bool {
	if !generateLock && !checkLock && !lockDiff && !lockClean && !lockUpdate {
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

	lockPath := "FluxFile.lock"

	if generateLock {
		lockFile, err := lock.GenerateWithPath(fluxFile, path, version)
		if err != nil {
			fmt.Printf("[ERROR] Failed to generate lock: %s\n", err.Error())
			return true
		}

		if err := lock.Save(lockFile, lockPath); err != nil {
			fmt.Printf("[ERROR] Failed to save lock: %s\n", err.Error())
			return true
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(lockFile, "", "  ")
			fmt.Println(string(data))
			return true
		}

		fmt.Printf("[✓] Lock file generated: %s (v%s)\n", lockPath, lockFile.Version)
		fmt.Printf("    Generated: %s\n", lockFile.Generated.Format("2006-01-02 15:04:05"))
		fmt.Printf("    OS/Arch: %s/%s\n", lockFile.Metadata.OS, lockFile.Metadata.Arch)
		fmt.Printf("    Tasks locked: %d\n", len(lockFile.Tasks))

		var totalInputs, totalOutputs int
		for taskName, taskLock := range lockFile.Tasks {
			totalInputs += len(taskLock.Inputs)
			totalOutputs += len(taskLock.Outputs)
			fmt.Printf("    - %s (%d inputs, %d outputs)\n",
				taskName, len(taskLock.Inputs), len(taskLock.Outputs))
		}
		fmt.Printf("    Total: %d inputs, %d outputs tracked\n", totalInputs, totalOutputs)
		return true
	}

	if lockUpdate {
		if updateTask == "" {
			fmt.Printf("[ERROR] Task name required for --lock-update\n")
			return true
		}

		lockFile, err := lock.Load(lockPath)
		if err != nil {
			fmt.Printf("[ERROR] Failed to load lock: %s\n", err.Error())
			return true
		}

		if err := lock.UpdateTask(lockFile, fluxFile, updateTask); err != nil {
			fmt.Printf("[ERROR] %s\n", err.Error())
			return true
		}

		if err := lock.Save(lockFile, lockPath); err != nil {
			fmt.Printf("[ERROR] Failed to save lock: %s\n", err.Error())
			return true
		}

		taskLock := lockFile.Tasks[updateTask]
		fmt.Printf("[✓] Updated task '%s' in lock file\n", updateTask)
		fmt.Printf("    Inputs: %d files\n", len(taskLock.Inputs))
		fmt.Printf("    Outputs: %d files\n", len(taskLock.Outputs))
		fmt.Printf("    Config hash: %s\n", taskLock.ConfigHash[:12])
		fmt.Printf("    Command hash: %s\n", taskLock.CommandHash[:12])
		return true
	}

	if lockDiff {
		lockFile, err := lock.Load(lockPath)
		if err != nil {
			fmt.Printf("[ERROR] Failed to load lock: %s\n", err.Error())
			return true
		}

		diffs := lock.ComputeDiff(lockFile, fluxFile)
		if len(diffs) == 0 {
			fmt.Println("[✓] No differences detected")
			return true
		}

		if jsonOutput {
			data, _ := json.MarshalIndent(diffs, "", "  ")
			fmt.Println(string(data))
			return true
		}

		fmt.Printf("[!] Found differences in %d task(s):\n\n", len(diffs))
		for _, diff := range diffs {
			fmt.Printf("Task: %s\n", diff.TaskName)

			if diff.ConfigChanged {
				fmt.Println("  [~] Task configuration changed")
			}
			if diff.CommandChanged {
				fmt.Println("  [~] Run commands changed")
			}

			for _, change := range diff.InputChanges {
				symbol := "~"
				if change.ChangeType == "missing" {
					symbol = "-"
				}
				fmt.Printf("  [%s] Input: %s (%s)\n", symbol, change.Path, change.ChangeType)
				if change.ChangeType == "size_changed" {
					fmt.Printf("      Size: %d -> %d bytes\n", change.OldSize, change.NewSize)
				}
			}

			for _, change := range diff.OutputChanges {
				symbol := "~"
				if change.ChangeType == "missing" {
					symbol = "-"
				}
				fmt.Printf("  [%s] Output: %s (%s)\n", symbol, change.Path, change.ChangeType)
				if change.ChangeType == "size_changed" {
					fmt.Printf("      Size: %d -> %d bytes\n", change.OldSize, change.NewSize)
				}
			}
			fmt.Println()
		}
		return true
	}

	if lockClean {
		lockFile, err := lock.Load(lockPath)
		if err != nil {
			fmt.Printf("[ERROR] Failed to load lock: %s\n", err.Error())
			return true
		}

		removed := lock.Clean(lockFile, fluxFile)
		if removed == 0 {
			fmt.Println("[✓] No stale tasks to clean")
			return true
		}

		if err := lock.Save(lockFile, lockPath); err != nil {
			fmt.Printf("[ERROR] Failed to save lock: %s\n", err.Error())
			return true
		}

		fmt.Printf("[✓] Removed %d stale task(s) from lock file\n", removed)
		return true
	}

	if checkLock {
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

		if jsonOutput {
			data, _ := json.MarshalIndent(changes, "", "  ")
			fmt.Println(string(data))
			return true
		}

		fmt.Printf("[⚠] Lock file verification failed - %d task(s) changed:\n", len(changes))
		for taskName, taskChanges := range changes {
			fmt.Printf("\n  Task: %s\n", taskName)
			for _, change := range taskChanges {
				fmt.Printf("    - %s\n", change)
			}
		}
		fmt.Println("\nRun 'flux --lock-diff' for detailed differences")
		return true
	}

	return false
}
