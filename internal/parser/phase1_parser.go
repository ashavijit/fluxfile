package parser

import (
	"strings"

	"github.com/ashavijit/fluxfile/internal/lexer"
)


func (p *Parser) parseDesc() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after desc")
		return ""
	}

	p.nextToken()

	if p.currentToken.Type == lexer.STRING {
		val := p.currentToken.Literal
		p.nextToken()
		return val
	} else if p.currentToken.Type == lexer.IDENT {
		var parts []string
		for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
			parts = append(parts, p.currentToken.Literal)
			p.nextToken()
		}
		return strings.Join(parts, " ")
	}

	return ""
}

func (p *Parser) parseParallel() bool {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after parallel")
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

func (p *Parser) parseIf() string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after if")
		return ""
	}

	p.nextToken()

	var condition strings.Builder
	for p.currentToken.Type != lexer.NEWLINE && p.currentToken.Type != lexer.EOF {
		if p.currentToken.Type == lexer.STRING {
			condition.WriteString("\"")
			condition.WriteString(p.currentToken.Literal)
			condition.WriteString("\"")
		} else {
			condition.WriteString(p.currentToken.Literal)
		}
		
		if p.peekToken.Type != lexer.NEWLINE && p.peekToken.Type != lexer.EOF {
			condition.WriteString(" ")
		}
		p.nextToken()
	}

	return strings.TrimSpace(condition.String())
}

func (p *Parser) parseCache() bool {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after cache")
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

func (p *Parser) parseInputs() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after inputs")
		return []string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var inputs []string

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.ENV || p.currentToken.Type == lexer.WATCH ||
			p.currentToken.Type == lexer.MATRIX || p.currentToken.Type == lexer.DOCKER ||
			p.currentToken.Type == lexer.REMOTE || p.currentToken.Type == lexer.OUTPUTS ||
			p.currentToken.Type == lexer.CACHE {
			break
		}

		if p.currentToken.Type == lexer.STRING || p.currentToken.Type == lexer.IDENT {
			inputs = append(inputs, p.currentToken.Literal)
			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return inputs
}

func (p *Parser) parseOutputs() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after outputs")
		return []string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var outputs []string

	for p.currentToken.Type != lexer.DEDENT && p.currentToken.Type != lexer.EOF {
		p.skipNewlines()

		if p.currentToken.Type == lexer.DEDENT || p.currentToken.Type == lexer.EOF {
			break
		}

		if p.currentToken.Type == lexer.DEPS || p.currentToken.Type == lexer.RUN ||
			p.currentToken.Type == lexer.ENV || p.currentToken.Type == lexer.WATCH ||
			p.currentToken.Type == lexer.MATRIX || p.currentToken.Type == lexer.DOCKER ||
			p.currentToken.Type == lexer.REMOTE || p.currentToken.Type == lexer.INPUTS ||
			p.currentToken.Type == lexer.CACHE {
			break
		}

		if p.currentToken.Type == lexer.STRING || p.currentToken.Type == lexer.IDENT {
			outputs = append(outputs, p.currentToken.Literal)
			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return outputs
}

func (p *Parser) parseWatchIgnore() []string {
	p.nextToken()

	if p.currentToken.Type != lexer.COLON {
		p.addError("expected : after ignore")
		return []string{}
	}

	p.nextToken()
	p.skipNewlines()

	if p.currentToken.Type == lexer.INDENT {
		p.nextToken()
	}

	var patterns []string

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
			patterns = append(patterns, p.currentToken.Literal)
			p.nextToken()
		} else {
			p.nextToken()
		}
	}

	if p.currentToken.Type == lexer.DEDENT {
		p.nextToken()
	}

	return patterns
}
