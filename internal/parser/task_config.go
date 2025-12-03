package parser

import (
	"strconv"

	"github.com/ashavijit/fluxfile/internal/ast"
	"github.com/ashavijit/fluxfile/internal/lexer"
)

func (p *Parser) parseProfileTask() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after profile_task")
		return ""
	}

	p.nextToken()

	switch p.currentToken.Type {
	case lexer.STRING, lexer.IDENT:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	}

	return ""
}

func (p *Parser) parseSecrets() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after secrets")
		return []string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var secrets []string

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.ENV || p.currentToken.Type == lexer.WATCH ||
			p.currentToken.Type == lexer.MATRIX || p.currentToken.Type == lexer.DOCKER ||
			p.currentToken.Type == lexer.REMOTE {
			break
		}

		if p.currentToken.Type == lexer.STRING || p.currentToken.Type == lexer.IDENT {
			secrets = append(secrets, p.currentToken.Literal)
			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return secrets
}

func (p *Parser) parsePre() []ast.Precondition {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after pre")
		return []ast.Precondition{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var preconditions []ast.Precondition

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.ENV {
			break
		}

		if p.currentToken.Type == lexer.IDENT {
			precType := p.currentToken.Literal
			p.nextToken()

			if p.currentToken.Type == lexer.COLON {
				p.nextToken()

				var value string
				switch p.currentToken.Type {
				case lexer.STRING, lexer.IDENT:
					value = p.currentToken.Literal
					p.nextToken()
				}

				preconditions = append(preconditions, ast.Precondition{
					Type:  precType,
					Value: value,
				})
			}
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return preconditions
}

func (p *Parser) parseRetries() int {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after retries")
		return 0
	}

	p.nextToken()

	switch p.currentToken.Type {
	case lexer.NUMBER, lexer.IDENT:
		val, err := strconv.Atoi(p.currentToken.Literal)
		if err != nil {
			p.addError("invalid number for retries")
			return 0
		}
		p.nextToken()
		return val
	}

	return 0
}

func (p *Parser) parseRetryDelay() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after retry_delay")
		return ""
	}

	p.nextToken()

	switch p.currentToken.Type {
	case lexer.STRING, lexer.IDENT:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	}

	return ""
}

func (p *Parser) parseTimeout() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after timeout")
		return ""
	}

	p.nextToken()

	switch p.currentToken.Type {
	case lexer.STRING, lexer.IDENT:
		val := p.currentToken.Literal
		p.nextToken()
		return val
	}

	return ""
}
