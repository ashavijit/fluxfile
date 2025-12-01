package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/ashavijit/fluxfile/internal/config"
	"github.com/ashavijit/fluxfile/internal/executor"
	"github.com/ashavijit/fluxfile/internal/logger"
	"github.com/ashavijit/fluxfile/internal/watcher"
)

var (
	version = "1.0.0"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

func main() {
	log := logger.New()

	taskName := flag.String("t", "", "Task to execute")
	profile := flag.String("p", "", "Profile to apply")
	listTasks := flag.Bool("l", false, "List all tasks")
	showTasks := flag.Bool("show", false, "Show all tasks with enhanced UI")
	watch := flag.Bool("w", false, "Watch mode")
	noCache := flag.Bool("no-cache", false, "Disable caching")
	fluxFilePath := flag.String("f", "", "Path to FluxFile")
	showVersion := flag.Bool("v", false, "Show version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Flux version %s\n", version)
		return
	}

	var path string
	var err error

	if *fluxFilePath != "" {
		path = *fluxFilePath
	} else {
		path, err = config.FindFluxFile()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	fluxFile, err := config.Load(path)
	if err != nil {
		log.Fatal(err.Error())
	}

	cacheDir := filepath.Join(".flux", "cache")
	exec, err := executor.New(fluxFile, cacheDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *showTasks || (len(flag.Args()) > 0 && flag.Args()[0] == "show") {
		showTasksEnhanced(exec)
		return
	}

	if *listTasks {
		tasks := exec.ListTasks()
		fmt.Println("Available tasks:")
		for _, task := range tasks {
			taskInfo, _ := exec.GetTaskInfo(task)
			if taskInfo != nil && taskInfo.Desc != "" {
				fmt.Printf("  %-20s %s\n", task, taskInfo.Desc)
			} else {
				fmt.Printf("  - %s\n", task)
			}
		}
		return
	}

	if *taskName == "" {
		if len(flag.Args()) > 0 {
			*taskName = flag.Args()[0]
		} else {
			log.Fatal("No task specified. Use -t <task> or provide task name as argument")
		}
	}

	task, err := exec.GetTaskInfo(*taskName)
	if err != nil {
		log.Fatal(err.Error())
	}

	if *watch && len(task.Watch) > 0 {
		log.Info(fmt.Sprintf("Starting watch mode for task: %s", *taskName))

		callback := func() {
			if err := exec.Execute(*taskName, *profile, !*noCache); err != nil {
				log.Error(err.Error())
			}
		}

		if err := exec.Execute(*taskName, *profile, !*noCache); err != nil {
			log.Error(err.Error())
		}

		w, err := watcher.New(task.Watch, callback)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := w.Start(); err != nil {
			log.Fatal(err.Error())
		}
	} else {
		if err := exec.Execute(*taskName, *profile, !*noCache); err != nil {
			log.Fatal(err.Error())
		}
	}
}
