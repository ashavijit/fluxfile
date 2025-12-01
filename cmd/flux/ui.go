package main

import (
	"fmt"
	"strings"

	"github.com/ashavijit/fluxfile/internal/executor"
)

func showTasksEnhanced(exec *executor.Executor) {
	tasks := exec.ListTasks()

	fmt.Println()
	fmt.Printf("%s╔════════════════════════════════════════════════════════════════╗%s\n", colorCyan, colorReset)
	fmt.Printf("%s║                        AVAILABLE TASKS                         ║%s\n", colorCyan, colorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════╝%s\n", colorCyan, colorReset)
	fmt.Println()

	maxNameLen := 0
	for _, taskName := range tasks {
		if len(taskName) > maxNameLen {
			maxNameLen = len(taskName)
		}
	}
	if maxNameLen < 15 {
		maxNameLen = 15
	}

	fmt.Printf("  %s%-*s  %s%-50s%s\n", colorYellow, maxNameLen, "TASK", colorGray, "DESCRIPTION", colorReset)
	fmt.Printf("  %s%s%s\n", colorGray, strings.Repeat("─", maxNameLen+54), colorReset)

	for _, taskName := range tasks {
		taskInfo, err := exec.GetTaskInfo(taskName)
		if err != nil {
			continue
		}

		desc := taskInfo.Desc
		if desc == "" {
			desc = colorGray + "(no description)" + colorReset
		}

		features := []string{}
		if taskInfo.Parallel {
			features = append(features, colorBlue+"⚡parallel"+colorReset)
		}
		if taskInfo.If != "" {
			features = append(features, colorYellow+"⚙ conditional"+colorReset)
		}
		if taskInfo.Cache {
			features = append(features, colorGreen+"💾 cached"+colorReset)
		}
		if taskInfo.Timeout != "" {
			features = append(features, colorYellow+"⏱ timeout"+colorReset)
		}
		if taskInfo.Retries > 0 {
			features = append(features, colorYellow+"🔄 retry"+colorReset)
		}
		if len(taskInfo.Deps) > 0 {
			features = append(features, fmt.Sprintf(colorGray+"→ %d deps"+colorReset, len(taskInfo.Deps)))
		}

		featureStr := ""
		if len(features) > 0 {
			featureStr = " " + colorGray + "[" + colorReset + strings.Join(features, " ") + colorGray + "]" + colorReset
		}

		fmt.Printf("  %s%-*s%s  %s%s\n", colorGreen, maxNameLen, taskName, colorReset, desc, featureStr)
	}

	fmt.Println()
	fmt.Printf("  %sTotal: %d tasks%s\n", colorGray, len(tasks), colorReset)
	fmt.Println()
	fmt.Printf("  %sRun a task:%s flux %s<task>%s\n", colorGray, colorReset, colorCyan, colorReset)
	fmt.Println()
}
