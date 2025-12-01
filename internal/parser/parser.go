package parser

import (
	"fmt"
	"strings"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/lexer"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken lexer.Token
	peekToken    lexer.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) skipNewlines() {
	for p.currentToken.Type == lexer.NEWLINE {
		p.nextToken()
	}
}

func (p *Parser) expectNewline() {
	if p.currentToken.Type == lexer.NEWLINE {
		p.nextToken()
	}
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, fmt.Sprintf("%s at line %d", msg, p.currentToken.Line))
}

func (p *Parser) Parse() (*ast.FluxFile, error) {
	fluxFile := ast.NewFluxFile()

	for p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT {
			p.nextToken()
			continue
		}

		switch p.currentToken.Type {
		case lexer.VAR:
			name, value := p.parseVarDecl()
			if name != "" {
				fluxFile.Vars[name] = value
			}
		case lexer.TASK:
			task := p.parseTask()
			if task.Name != "" {
				fluxFile.Tasks = append(fluxFile.Tasks, task)
			}
		case lexer.PROFILE:
			profile := p.parseProfile()
			if profile.Name != "" {
				fluxFile.Profiles = append(fluxFile.Profiles, profile)
			}
		case lexer.INCLUDE:
			include := p.parseInclude()
			if include != "" {
				fluxFile.Includes = append(fluxFile.Includes, include)
			}
		case lexer.EOF:
			break
		default:
			p.nextToken()
		}
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("parse errors: %s", strings.Join(p.errors, "; "))
	}

	return fluxFile, nil
}

func (p *Parser) parseVarDecl() (string, string) {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("expected identifier after var, got %s", p.currentToken.Type))
		return "", ""
	}

	name := p.currentToken.Literal
	p.nextToken()

	if p.currentToken.Type != lexer.EQUALS {
		p.addError(fmt.Sprintf("expected =, got %s", p.currentToken.Type))
		return "", ""
	}

	p.nextToken()
	value := p.parseExpr()

	return name, value
}

func (p *Parser) parseExpr() string {
	switch p.currentToken.Type {
	case lexer.STRING:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	case lexer.NUMBER:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	case lexer.IDENT:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	case lexer.DOLLAR:
		return p.parseShellExpr()
	default:
		p.addError(fmt.Sprintf("unexpected expression token %s", p.currentToken.Type))
		p.nextToken()
		return ""
	}
}

func (p *Parser) parseShellExpr() string {
	if p.currentToken.Type != lexer.DOLLAR {
		return ""
	}

	p.nextToken()

	if p.currentToken.Type != lexer.LPAREN {
		p.addError("expected ( after $")
		return ""
	}

	p.nextToken()

	// Skip whitespace if present
	p.skipNewlines()
	for p.currentToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	if p.currentToken.Type != lexer.SHELL && p.currentToken.Type != lexer.IDENT {
		p.addError("expected shell keyword")
		return ""
	}

	// Verify it's actually "shell"
	if p.currentToken.Type == lexer.IDENT && p.currentToken.Literal != "shell" {
		p.addError("expected shell keyword")
		return ""
	}

	p.nextToken()

	// Skip whitespace (spec: SP)
	if p.currentToken.Type == lexer.NEWLINE {
		p.nextToken()
	}

	if p.currentToken.Type != lexer.STRING {
		p.addError("expected string after shell")
		return ""
	}

	command := p.currentToken.Literal
	p.nextToken()

	if p.currentToken.Type != lexer.RPAREN {
		p.addError("expected ) after shell command")
		return ""
	}

	p.nextToken()

	return fmt.Sprintf("$(shell %q)", command)
}

func (p *Parser) parseTask() ast.Task {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("expected task name, got %s", p.currentToken.Type))
		return ast.Task{}
	}

	task := ast.NewTask(p.currentToken.Literal)
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError(fmt.Sprintf("expected :, got %s", p.currentToken.Type))
		return ast.Task{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	p.parseTaskBody(&task)

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return task
}

func (p *Parser) parseTaskBody(task *ast.Task) {
	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.TASK || p.currentToken.Type == lexer.PROFILE ||
			p.currentToken.Type == lexer.VAR || p.currentToken.Type == lexer.INCLUDE {
			break
		}

				switch p.currentToken.Type {
	case lexer.DESC:
		task.Desc = p.parseDesc()
	case lexer.DEPS:
		task.Deps = p.parseDeps()
	case lexer.PARALLEL:
		task.Parallel = p.parseParallel()
	case lexer.IF:
		task.If = p.parseIf()
	case lexer.RUN:
		task.Run = p.parseRun()
	case lexer.ENV:
		task.Env = p.parseEnv()
	case lexer.WATCH:
		task.Watch = p.parseWatch()
	case lexer.IGNORE:
		task.WatchIgnore = p.parseWatchIgnore()
	case lexer.MATRIX:
		task.Matrix = p.parseMatrix()
	case lexer.CACHE:
		task.Cache = p.parseCache()
	case lexer.INPUTS:
		task.Inputs = p.parseInputs()
	case lexer.OUTPUTS:
		task.Outputs = p.parseOutputs()
	case lexer.PROFILE_TASK:
		task.Profile = p.parseProfileTask()
	case lexer.SECRETS:
		task.Secrets = p.parseSecrets()
	case lexer.PRE:
		task.Pre = p.parsePre()
	case lexer.RETRIES:
		task.Retries = p.parseRetries()
	case lexer.RETRY_DELAY:
		task.RetryDelay = p.parseRetryDelay()
	case lexer.TIMEOUT:
		task.Timeout = p.parseTimeout()
	case lexer.DOCKER:
		task.Docker = p.parseDocker()
	case lexer.REMOTE:
		task.Remote = p.parseRemote()
	default:
		p.nextToken()
	}


	}
}

func (p *Parser) parseDeps() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after deps")
		return []string{}
	}

	p.nextToken()

	var deps []string

	for {
		if p.currentToken.Type != lexer.IDENT {
			break
		}

		deps = append(deps, p.currentToken.Literal)
		p.nextToken()

		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	return deps
}

func (p *Parser) parseRun() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after run")
		return []string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var commands []string

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		command := p.parseCommand()
		if command != "" {
			commands = append(commands, command)
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return commands
}

func (p *Parser) parseCommand() string {
	var sb strings.Builder

	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF && p.currentToken.Type != lexer.DEDENT {
		text := p.currentToken.Literal
		length := len(text)

		if p.currentToken.Type == lexer.STRING {
			text = "\"" + text + "\""
			length += 2 // Account for quotes
		}

		sb.WriteString(text)

		// Check if we need to add a space before the next token
		// We use the column information from the lexer to determine if there was a space in the source
		currentEnd := p.currentToken.Column + length
		nextStart := p.peekToken.Column

		// Only add space if the next token is on the same line and there is a gap
		if p.peekToken.Line == p.currentToken.Line && nextStart > currentEnd {
			sb.WriteString(" ")
		}

		p.nextToken()
	}

	return strings.TrimSpace(sb.String())
}

func (p *Parser) parseEnv() map[string]string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after env")
		return map[string]string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	env := make(map[string]string)

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.WATCH || p.currentToken.Type == lexer.MATRIX ||
			p.currentToken.Type == lexer.DOCKER || p.currentToken.Type == lexer.REMOTE {
			break
		}

		if p.currentToken.Type != lexer.IDENT {
			p.nextToken()
			continue
		}

		key := p.currentToken.Literal
		p.nextToken()

		if p.currentToken.Type != lexer.EQUALS {
			p.addError("expected = in env block")
			continue
		}

		p.nextToken()
		value := p.parseExpr()

		env[key] = value
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return env
}

func (p *Parser) parseWatch() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after watch")
		return []string{}
	}

	p.nextToken()

	var patterns []string

	if p.currentToken.Type == lexer.STRING {
		patterns = append(patterns, p.currentToken.Literal)
		p.nextToken()
	} else if p.currentToken.Type == lexer.IDENT {
		patterns = append(patterns, p.currentToken.Literal)
		p.nextToken()
	}

	return patterns
}

func (p *Parser) parseMatrix() *ast.Matrix {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after matrix")
		return nil
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	matrix := ast.NewMatrix()

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.ENV || p.currentToken.Type == lexer.WATCH ||
			p.currentToken.Type == lexer.DOCKER || p.currentToken.Type == lexer.REMOTE {
			break
		}

		if p.currentToken.Type != lexer.IDENT {
			p.nextToken()
			continue
		}

		key := p.currentToken.Literal
		p.nextToken()

		if p.currentToken.Type != lexer.COLON {
			p.addError("expected : in matrix block")
			continue
		}

		p.nextToken()

		var values []string
		for {
			if p.currentToken.Type != lexer.IDENT && p.currentToken.Type != lexer.STRING {
				break
			}

			values = append(values, p.currentToken.Literal)
			p.nextToken()

			if p.currentToken.Type == lexer.COMMA {
				p.nextToken()
			} else {
				break
			}
		}

		matrix.Dimensions[key] = values
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return matrix
}

func (p *Parser) parseDocker() bool {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after docker")
		return false
	}

	p.nextToken()

	if p.currentToken.Type == lexer.IDENT {
		val := p.currentToken.Literal == "true"
		p.nextToken()
		return val
	}

	return false
}

func (p *Parser) parseRemote() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after remote")
		return ""
	}

	p.nextToken()

	if p.currentToken.Type != lexer.STRING {
		p.addError("expected string after remote:")
		return ""
	}

	remote := p.currentToken.Literal
	p.nextToken()

	return remote
}

func (p *Parser) parseProfile() ast.Profile {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		p.addError(fmt.Sprintf("expected profile name, got %s", p.currentToken.Type))
		return ast.Profile{}
	}

	profile := ast.NewProfile(p.currentToken.Literal)
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError(fmt.Sprintf("expected :, got %s", p.currentToken.Type))
		return ast.Profile{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.ENV {
			profile.Env = p.parseEnv()
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return profile
}

func (p *Parser) parseInclude() string {
	p.nextToken()

	if p.currentToken.Type != lexer.STRING {
		p.addError(fmt.Sprintf("expected string after include, got %s", p.currentToken.Type))
		return ""
	}

	include := p.currentToken.Literal
	p.nextToken()

	return include
}


