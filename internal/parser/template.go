package parser

import (
	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/lexer"
)

// parseTemplate parses a template definition block.
// Templates are reusable task definitions that can be extended by tasks.
//
// Syntax:
//
//	template name:
//	    cache: true
//	    inputs:
//	        **/*.go
func (p *Parser) parseTemplate() ast.Template {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		p.addError("expected template name")
		return ast.Template{}
	}

	template := ast.NewTemplate(p.currentToken.Literal)
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after template name")
		return ast.Template{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	p.parseTemplateBody(&template)

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return template
}

func (p *Parser) parseTemplateBody(template *ast.Template) {
	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.TASK || p.currentToken.Type == lexer.PROFILE ||
			p.currentToken.Type == lexer.VAR || p.currentToken.Type == lexer.INCLUDE ||
			p.currentToken.Type == lexer.TEMPLATE || p.currentToken.Type == lexer.GROUP {
			break
		}

		switch p.currentToken.Type {
		case lexer.DESC:
			template.Desc = p.parseDesc()
		case lexer.DEPS:
			template.Deps = p.parseDeps()
		case lexer.PARALLEL:
			template.Parallel = p.parseParallel()
		case lexer.ENV:
			template.Env = p.parseEnv()
		case lexer.CACHE:
			template.Cache = p.parseCache()
		case lexer.INPUTS:
			template.Inputs = p.parseInputs()
		case lexer.OUTPUTS:
			template.Outputs = p.parseOutputs()
		case lexer.SECRETS:
			template.Secrets = p.parseSecrets()
		case lexer.PRE:
			template.Pre = p.parsePre()
		case lexer.RETRIES:
			template.Retries = p.parseRetries()
		case lexer.RETRY_DELAY:
			template.RetryDelay = p.parseRetryDelay()
		case lexer.TIMEOUT:
			template.Timeout = p.parseTimeout()
		case lexer.DOCKER:
			template.Docker = p.parseDocker()
		case lexer.REMOTE:
			template.Remote = p.parseRemote()
		case lexer.BEFORE:
			template.Before = p.parseBefore()
		case lexer.AFTER:
			template.After = p.parseAfter()
		default:
			p.nextToken()
		}
	}
}

// parseTaskGroup parses a group definition block.
// Groups organize related tasks under a namespace.
//
// Syntax:
//
//	group frontend:
//	    tasks: install, build, test
func (p *Parser) parseTaskGroup() ast.TaskGroup {
	p.nextToken()

	if p.currentToken.Type != lexer.IDENT {
		p.addError("expected group name")
		return ast.TaskGroup{}
	}

	group := ast.NewTaskGroup(p.currentToken.Literal)
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after group name")
		return ast.TaskGroup{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	p.parseGroupBody(&group)

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return group
}

// parseGroupBody parses the body of a group definition
func (p *Parser) parseGroupBody(group *ast.TaskGroup) {
	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.TASK || p.currentToken.Type == lexer.PROFILE ||
			p.currentToken.Type == lexer.VAR || p.currentToken.Type == lexer.INCLUDE ||
			p.currentToken.Type == lexer.TEMPLATE || p.currentToken.Type == lexer.GROUP {
			break
		}

		if p.currentToken.Type == lexer.IDENT && p.currentToken.Literal == "tasks" {
			p.nextToken()
			if p.currentToken.Type == lexer.COLON {
				p.nextToken()
				group.Tasks = p.parseGroupTasks()
			}
		} else {
			p.nextToken()
		}
	}
}

// parseGroupTasks parses the tasks list in a group definition
func (p *Parser) parseGroupTasks() []string {
	var tasks []string

	for {
		if p.currentToken.Type != lexer.IDENT {
			break
		}

		tasks = append(tasks, p.currentToken.Literal)
		p.nextToken()

		if p.currentToken.Type == lexer.COMMA {
			p.nextToken()
		} else {
			break
		}
	}

	return tasks
}
