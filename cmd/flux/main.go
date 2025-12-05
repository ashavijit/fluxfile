package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/ashavijit/fluxfile/internal/config"
	"github.com/ashavijit/fluxfile/internal/executor"
	fluxinit "github.com/ashavijit/fluxfile/internal/init"
	"github.com/ashavijit/fluxfile/internal/logger"
	"github.com/ashavijit/fluxfile/internal/report"
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
	colorRed    = "\033[31m"
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
	generateLock := flag.Bool("lock", false, "Generate dependency lock file")
	checkLock := flag.Bool("check-lock", false, "Verify lock file")
	lockUpdate := flag.Bool("lock-update", false, "Update specific task in lock file")
	updateTask := flag.String("task", "", "Task name for --lock-update")
	lockDiff := flag.Bool("lock-diff", false, "Show detailed diff between lock and current state")
	lockClean := flag.Bool("lock-clean", false, "Remove stale tasks from lock file")
	jsonOutput := flag.Bool("json", false, "Output in JSON format")
	runTUI := flag.Bool("tui", false, "Run interactive TUI mode")
	dryRun := flag.Bool("dry-run", false, "Simulate task execution")
	completion := flag.String("completion", "", "Generate shell completion script (bash, zsh, fish, powershell)")
	initCmd := flag.Bool("init", false, "Initialize a new FluxFile")
	initTemplate := flag.String("template", "", "Template for init (go, node, python, rust, generic)")
	showReport := flag.Bool("report", false, "Show execution report after task completion")
	reportJSON := flag.String("report-json", "", "Save execution report as JSON to specified path")

	flag.Parse()

	if *showVersion {
		fmt.Printf("Flux version %s\n", version)
		return
	}

	if *completion != "" {
		generateCompletion(*completion)
		return
	}

	if *initCmd || (len(flag.Args()) > 0 && flag.Args()[0] == "init") {
		cfg := fluxinit.Config{
			Template:  *initTemplate,
			Directory: ".",
		}
		if err := fluxinit.Run(cfg); err != nil {
			log.Fatal(err.Error())
		}
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
	exec, err := executor.New(fluxFile, cacheDir, *dryRun)
	if err != nil {
		log.Fatal(err.Error())
	}

	if handleLockCommands(*generateLock, *checkLock, *lockUpdate, *lockDiff, *lockClean, *updateTask, *fluxFilePath, *jsonOutput) {
		return
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

	if *runTUI {
		runInteractiveTUI(exec, *taskName, *profile, !*noCache)
		return
	}

	var collector *report.Collector
	if *showReport || *reportJSON != "" {
		collector = report.NewCollector()
		exec.SetCollector(collector)
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

	if collector != nil {
		rep := collector.Generate()
		if *showReport {
			rep.Print()
		}
		if *reportJSON != "" {
			if err := rep.WriteJSON(*reportJSON); err != nil {
				log.Error(fmt.Sprintf("Failed to write report JSON: %s", err.Error()))
			} else {
				fmt.Printf("Report saved to %s\n", *reportJSON)
			}
		}
	}
}

func generateCompletion(shell string) {
	switch shell {
	case "bash":
		fmt.Println(`_flux_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="$(flux -l | grep '  -' | awk '{print $2}')"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "-t -p -l -show -w -no-cache -f -v -lock -check-lock -lock-update -lock-diff -lock-clean -json -tui -dry-run" -- ${cur}) )
        return 0
    fi

    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
}
complete -F _flux_completion flux`)
	case "zsh":
		fmt.Println(`#compdef flux
_flux() {
    local -a tasks
    tasks=("${(@f)$(flux -l | grep '  -' | awk '{print $2}')}")
    _arguments \
        '-t[Task to execute]' \
        '-p[Profile to apply]' \
        '-l[List all tasks]' \
        '-show[Show all tasks with enhanced UI]' \
        '-w[Watch mode]' \
        '-no-cache[Disable caching]' \
        '-f[Path to FluxFile]' \
        '-v[Show version]' \
        '-lock[Generate dependency lock file]' \
        '-check-lock[Verify lock file]' \
        '-lock-update[Update specific task in lock file]' \
        '-lock-diff[Show detailed diff between lock and current state]' \
        '-lock-clean[Remove stale tasks from lock file]' \
        '-json[Output in JSON format]' \
        '-tui[Run interactive TUI mode]' \
        '-dry-run[Simulate task execution]' \
        '1: :($tasks)'
}
_flux`)
	case "fish":
		fmt.Println(`function __fish_flux_tasks
    flux -l | grep '  -' | awk '{print $2}'
end

complete -f -c flux -a "(__fish_flux_tasks)"
complete -c flux -s t -d "Task to execute"
complete -c flux -s p -d "Profile to apply"
complete -c flux -s l -d "List all tasks"
complete -c flux -s show -d "Show all tasks with enhanced UI"
complete -c flux -s w -d "Watch mode"
complete -c flux -s no-cache -d "Disable caching"
complete -c flux -s f -d "Path to FluxFile"
complete -c flux -s v -d "Show version"
complete -c flux -s lock -d "Generate dependency lock file"
complete -c flux -s check-lock -d "Verify lock file"
complete -c flux -s lock-update -d "Update specific task in lock file"
complete -c flux -s lock-diff -d "Show detailed diff between lock and current state"
complete -c flux -s lock-clean -d "Remove stale tasks from lock file"
complete -c flux -s json -d "Output in JSON format"
complete -c flux -s tui -d "Run interactive TUI mode"
complete -c flux -s dry-run -d "Simulate task execution"`)
	case "powershell":
		fmt.Println(`Register-ArgumentCompleter -Native -CommandName flux -ScriptBlock {
    param($commandName, $wordToComplete, $cursorPosition)
    $tasks = flux -l | Select-String "  -" | ForEach-Object { $_.ToString().Trim().Substring(2) }
    $tasks | Where-Object { $_ -like "$wordToComplete*" }
}`)
	default:
		fmt.Printf("Unsupported shell: %s\n", shell)
	}
}
