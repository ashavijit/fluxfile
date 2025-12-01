# Enhanced Task Display

Use `flux show` or `flux -show` for a beautiful terminal UI showing all tasks.

## Usage

```bash
# Show all tasks with enhanced UI
flux show
flux -show

# With custom FluxFile
flux -f FluxFile.example show
```

## Features

- ✅ Colored, bordered display
- ✅ Task descriptions
- ✅ Feature badges:
  - ⚡ Parallel execution
  - ⚙️ Conditional execution
  - 💾 Cached tasks
  - ⏱️ Timeout enabled
  - 🔄 Retry enabled
  - → Dependency count
- ✅ Total task count
- ✅ Usage hint

## Example Output

```
╔════════════════════════════════════════════════════════════════╗
║                        AVAILABLE TASKS                         ║
╚════════════════════════════════════════════════════════════════╝

  TASK              DESCRIPTION
  ─────────────────────────────────────────────────────────────

  build             Build the Go binary [💾 cached → 2 deps]
  deploy            Deploy to production [⚙ conditional ⏱ timeout]
  ci                Run CI tasks [⚡parallel → 3 deps]

  Total: 11 tasks

  Run a task: flux <task>
```
