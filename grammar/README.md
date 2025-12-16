# FluxFile Language Support

This directory contains the TextMate grammar for FluxFile syntax highlighting.

## Supported Features

- **Keywords**: `task`, `var`, `profile`, `include`
- **Properties**: `desc`, `deps`, `run`, `env`, `inputs`, `outputs`, `cache`, `watch`, `docker`, `remote`, `matrix`, `parallel`, `if`, `ignore`
- **Variables**: `${VAR}` interpolation
- **Shell commands**: `$(shell "command")`
- **Comments**: `# comment`

## Installing Syntax Highlighting

### VS Code Extension

1. Install the FluxFile VS Code extension (coming soon)
2. Or manually copy the grammar to your VS Code extensions

### GitHub Linguist (Official Recognition)

To get GitHub to recognize FluxFile files natively:

1. **Fork** [github/linguist](https://github.com/github/linguist)
2. Add entry to `lib/linguist/languages.yml`:
   ```yaml
   FluxFile:
     type: data
     color: "#6366f1"
     extensions:
       - ".flux"
     filenames:
       - FluxFile
       - fluxfile
     tm_scope: source.fluxfile
     ace_mode: yaml
     codemirror_mode: yaml
     codemirror_mime_type: text/x-yaml
   ```
3. Add grammar to `vendor/grammars/`
4. Submit a PR

### Temporary: Use gitattributes

Add to your repo's `.gitattributes` to make GitHub highlight FluxFile as YAML:

```
FluxFile linguist-language=YAML
FluxFile.* linguist-language=YAML
*.flux linguist-language=YAML
```

## Grammar File

- `fluxfile.tmLanguage.json` - TextMate grammar (JSON format)

## Colors (Theme Reference)

| Scope | Purpose | Example |
|-------|---------|---------|
| `keyword.control.task` | Task keyword | `task` |
| `entity.name.function.task` | Task name | `build` |
| `keyword.other.var` | Variable keyword | `var` |
| `variable.other.constant` | Variable name | `PROJECT` |
| `keyword.other.property` | Properties | `deps`, `cache` |
| `keyword.control.section` | Sections | `run`, `env` |
| `string.unquoted.command` | Shell commands | `go build` |
| `variable.other.interpolation` | Interpolated vars | `${VAR}` |
