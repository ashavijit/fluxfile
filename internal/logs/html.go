package logs

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flux Logs</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
            background: #0d1117;
            color: #c9d1d9;
            padding: 20px;
            line-height: 1.6;
        }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 {
            color: #58a6ff;
            border-bottom: 1px solid #30363d;
            padding-bottom: 16px;
            margin-bottom: 24px;
        }
        .task-card {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 8px;
            margin-bottom: 16px;
            overflow: hidden;
        }
        .task-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 16px;
            background: #21262d;
            cursor: pointer;
        }
        .task-header:hover { background: #30363d; }
        .task-name { font-weight: bold; font-size: 16px; }
        .task-status {
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: bold;
        }
        .status-success { background: #238636; color: #fff; }
        .status-failed { background: #da3633; color: #fff; }
        .status-running { background: #bf8700; color: #fff; }
        .task-meta { color: #8b949e; font-size: 12px; margin-top: 4px; }
        .task-body { display: none; padding: 16px; border-top: 1px solid #30363d; }
        .task-card.expanded .task-body { display: block; }
        .log-entry {
            padding: 8px 12px;
            margin: 4px 0;
            border-radius: 4px;
            font-size: 13px;
            background: #0d1117;
        }
        .log-timestamp { color: #8b949e; margin-right: 12px; }
        .log-level {
            display: inline-block;
            width: 50px;
            font-weight: bold;
        }
        .level-info { color: #58a6ff; }
        .level-cmd { color: #a371f7; }
        .level-error { color: #f85149; }
        .level-warn { color: #d29922; }
        .log-message { color: #c9d1d9; }
        .log-command {
            background: #1f2428;
            padding: 8px;
            border-radius: 4px;
            margin-top: 4px;
            color: #7ee787;
        }
        .log-duration { color: #8b949e; font-size: 11px; }
        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #8b949e;
        }
        .summary {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 16px;
            margin-bottom: 24px;
        }
        .summary-card {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 8px;
            padding: 16px;
            text-align: center;
        }
        .summary-value { font-size: 32px; font-weight: bold; color: #58a6ff; }
        .summary-label { color: #8b949e; font-size: 12px; margin-top: 4px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Flux Execution Logs</h1>
        
        {{if .Tasks}}
        <div class="summary">
            <div class="summary-card">
                <div class="summary-value">{{.TotalTasks}}</div>
                <div class="summary-label">Total Tasks</div>
            </div>
            <div class="summary-card">
                <div class="summary-value" style="color: #238636">{{.SuccessCount}}</div>
                <div class="summary-label">Successful</div>
            </div>
            <div class="summary-card">
                <div class="summary-value" style="color: #da3633">{{.FailedCount}}</div>
                <div class="summary-label">Failed</div>
            </div>
        </div>
        
        {{range .Tasks}}
        <div class="task-card" onclick="this.classList.toggle('expanded')">
            <div class="task-header">
                <div>
                    <div class="task-name">{{.TaskName}}</div>
                    <div class="task-meta">
                        Started: {{.StartTime.Format "2006-01-02 15:04:05"}}
                        {{if not .EndTime.IsZero}}
                        | Duration: {{duration .StartTime .EndTime}}
                        {{end}}
                    </div>
                </div>
                <span class="task-status status-{{.Status}}">{{.Status}}</span>
            </div>
            <div class="task-body">
                {{if .Entries}}
                {{range .Entries}}
                <div class="log-entry">
                    <span class="log-timestamp">{{.Timestamp.Format "15:04:05.000"}}</span>
                    <span class="log-level level-{{.Level}}">{{.Level}}</span>
                    {{if .Command}}
                    <div class="log-command">$ {{.Command}}</div>
                    {{if .Duration}}<span class="log-duration">{{.Duration}}ms</span>{{end}}
                    {{else}}
                    <span class="log-message">{{.Message}}</span>
                    {{end}}
                </div>
                {{end}}
                {{else}}
                <div class="empty-state">No log entries</div>
                {{end}}
            </div>
        </div>
        {{end}}
        {{else}}
        <div class="empty-state">
            <h2>No logs found</h2>
            <p>Run some tasks to generate logs</p>
        </div>
        {{end}}
    </div>
    <script>
        // Auto-expand first task
        const first = document.querySelector('.task-card');
        if (first) first.classList.add('expanded');
    </script>
</body>
</html>`

type HTMLData struct {
	Tasks        []*TaskLog
	TotalTasks   int
	SuccessCount int
	FailedCount  int
}

func GenerateHTML(logs []*TaskLog) (string, error) {
	// Sort by start time descending
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].StartTime.After(logs[j].StartTime)
	})

	data := HTMLData{
		Tasks:      logs,
		TotalTasks: len(logs),
	}

	for _, log := range logs {
		switch log.Status {
		case "success":
			data.SuccessCount++
		case "failed":
			data.FailedCount++
		}
	}

	funcs := template.FuncMap{
		"duration": func(start, end time.Time) string {
			d := end.Sub(start)
			if d < time.Second {
				return fmt.Sprintf("%dms", d.Milliseconds())
			}
			return fmt.Sprintf("%.2fs", d.Seconds())
		},
	}

	tmpl, err := template.New("logs").Funcs(funcs).Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	path := filepath.Join(GetLogDir(), "logs.html")
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return "", err
	}

	return path, nil
}
