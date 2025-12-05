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
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #0d1117;
            color: #c9d1d9;
            padding: 24px;
            line-height: 1.5;
        }
        .container { max-width: 1400px; margin: 0 auto; }
        h1 {
            color: #58a6ff;
            font-size: 24px;
            margin-bottom: 24px;
            display: flex;
            align-items: center;
            gap: 12px;
        }
        h1::before { content: ""; display: inline-block; width: 8px; height: 8px; background: #238636; border-radius: 50%; }
        .summary {
            display: grid;
            grid-template-columns: repeat(4, 1fr);
            gap: 16px;
            margin-bottom: 32px;
        }
        .summary-card {
            background: #161b22;
            border: 1px solid #30363d;
            border-radius: 8px;
            padding: 20px;
        }
        .summary-value { font-size: 36px; font-weight: 600; }
        .summary-label { color: #8b949e; font-size: 13px; margin-top: 4px; }
        .success { color: #3fb950; }
        .failed { color: #f85149; }
        .primary { color: #58a6ff; }
        
        table {
            width: 100%;
            border-collapse: collapse;
            background: #161b22;
            border-radius: 8px;
            overflow: hidden;
            margin-bottom: 24px;
        }
        th {
            text-align: left;
            padding: 14px 16px;
            background: #21262d;
            color: #8b949e;
            font-weight: 500;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
            border-bottom: 1px solid #30363d;
        }
        td {
            padding: 14px 16px;
            border-bottom: 1px solid #21262d;
            font-size: 14px;
        }
        tr:hover { background: #1c2128; }
        tr:last-child td { border-bottom: none; }
        
        .status-badge {
            display: inline-block;
            padding: 4px 10px;
            border-radius: 16px;
            font-size: 12px;
            font-weight: 500;
        }
        .status-success { background: rgba(63, 185, 80, 0.15); color: #3fb950; }
        .status-failed { background: rgba(248, 81, 73, 0.15); color: #f85149; }
        .status-running { background: rgba(210, 153, 34, 0.15); color: #d29922; }
        
        .task-name { font-weight: 500; color: #f0f6fc; }
        .mono { font-family: 'SF Mono', 'Fira Code', monospace; font-size: 13px; }
        .muted { color: #8b949e; }
        
        .details-row { display: none; }
        .details-row.show { display: table-row; }
        .details-cell {
            padding: 0 !important;
            background: #0d1117;
        }
        .log-table {
            width: 100%;
            margin: 0;
            border-radius: 0;
        }
        .log-table th { background: #161b22; }
        .log-table td { padding: 10px 16px; font-size: 13px; }
        
        .level-info { color: #58a6ff; }
        .level-cmd { color: #a371f7; }
        .level-error { color: #f85149; }
        .level-warn { color: #d29922; }
        
        .cmd-text { color: #7ee787; font-family: 'SF Mono', monospace; }
        
        .toggle-btn {
            background: none;
            border: 1px solid #30363d;
            color: #8b949e;
            padding: 4px 10px;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
        }
        .toggle-btn:hover { background: #21262d; color: #c9d1d9; }
        
        .empty-state {
            text-align: center;
            padding: 80px 20px;
            color: #8b949e;
        }
        .empty-state h2 { color: #c9d1d9; margin-bottom: 8px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Flux Execution Logs</h1>
        
        {{if .Tasks}}
        <div class="summary">
            <div class="summary-card">
                <div class="summary-value primary">{{.TotalTasks}}</div>
                <div class="summary-label">Total Executions</div>
            </div>
            <div class="summary-card">
                <div class="summary-value success">{{.SuccessCount}}</div>
                <div class="summary-label">Successful</div>
            </div>
            <div class="summary-card">
                <div class="summary-value failed">{{.FailedCount}}</div>
                <div class="summary-label">Failed</div>
            </div>
            <div class="summary-card">
                <div class="summary-value primary">{{if .Tasks}}{{duration (index .Tasks 0).StartTime (index .Tasks 0).EndTime}}{{else}}-{{end}}</div>
                <div class="summary-label">Last Duration</div>
            </div>
        </div>
        
        <table>
            <thead>
                <tr>
                    <th style="width: 40px;"></th>
                    <th>Task</th>
                    <th>Status</th>
                    <th>Started</th>
                    <th>Duration</th>
                    <th>Entries</th>
                </tr>
            </thead>
            <tbody>
                {{range $i, $task := .Tasks}}
                <tr onclick="toggleDetails({{$i}})" style="cursor: pointer;">
                    <td><button class="toggle-btn" id="btn-{{$i}}">+</button></td>
                    <td class="task-name">{{.TaskName}}</td>
                    <td><span class="status-badge status-{{.Status}}">{{.Status}}</span></td>
                    <td class="mono muted">{{.StartTime.Format "Jan 02 15:04:05"}}</td>
                    <td class="mono">{{if not .EndTime.IsZero}}{{duration .StartTime .EndTime}}{{else}}-{{end}}</td>
                    <td class="muted">{{len .Entries}} entries</td>
                </tr>
                <tr class="details-row" id="details-{{$i}}">
                    <td colspan="6" class="details-cell">
                        <div style="padding: 12px 16px; background: #161b22; border-bottom: 1px solid #21262d;">
                            <span style="color: #8b949e; margin-right: 20px;">
                                <strong>Working Dir:</strong> {{if .WorkDir}}{{.WorkDir}}{{else}}./{{end}}
                            </span>
                            {{if .Profile}}<span style="color: #a371f7; margin-right: 20px;"><strong>Profile:</strong> {{.Profile}}</span>{{end}}
                            {{if .CacheHit}}<span style="color: #3fb950; margin-right: 20px;">Cache Hit</span>{{end}}
                            {{if .DepsCount}}<span style="color: #8b949e; margin-right: 20px;"><strong>Dependencies:</strong> {{.DepsCount}}</span>{{end}}
                            {{if .Error}}<span style="color: #f85149;"><strong>Error:</strong> {{.Error}}</span>{{end}}
                        </div>
                        {{if .Entries}}
                        <table class="log-table">
                            <thead>
                                <tr>
                                    <th style="width: 120px;">Time</th>
                                    <th style="width: 80px;">Level</th>
                                    <th>Message</th>
                                    <th style="width: 80px;">Exit</th>
                                    <th style="width: 100px;">Duration</th>
                                </tr>
                            </thead>
                            <tbody>
                                {{range .Entries}}
                                <tr>
                                    <td class="mono muted">{{.Timestamp.Format "15:04:05.000"}}</td>
                                    <td><span class="level-{{.Level}}">{{.Level}}</span></td>
                                    <td>{{if .Command}}<span class="cmd-text">$ {{.Command}}</span>{{if .Output}}<div style="color: #8b949e; font-size: 12px; margin-top: 4px; white-space: pre-wrap;">{{.Output}}</div>{{end}}{{else}}{{.Message}}{{end}}</td>
                                    <td class="mono {{if .ExitCode}}failed{{else}}muted{{end}}">{{if .Command}}{{.ExitCode}}{{else}}-{{end}}</td>
                                    <td class="mono muted">{{if .Duration}}{{.Duration}}ms{{else}}-{{end}}</td>
                                </tr>
                                {{end}}
                            </tbody>
                        </table>
                        {{else}}
                        <div style="padding: 20px; text-align: center; color: #8b949e;">No log entries</div>
                        {{end}}
                    </td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <div class="empty-state">
            <h2>No logs found</h2>
            <p>Run some tasks to generate logs</p>
        </div>
        {{end}}
    </div>
    <script>
        function toggleDetails(i) {
            const row = document.getElementById('details-' + i);
            const btn = document.getElementById('btn-' + i);
            row.classList.toggle('show');
            btn.textContent = row.classList.contains('show') ? '-' : '+';
        }
        // Auto-expand first
        if (document.querySelector('.details-row')) toggleDetails(0);
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
