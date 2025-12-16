package lexer

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input         string
	position      int
	readPosition  int
	ch            byte
	line          int
	column        int
	indentStack   []int
	pendingDedent int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:       input,
		line:        1,
		column:      0,
		indentStack: []int{0},
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) NextToken() Token {
	if l.pendingDedent > 0 {
		l.pendingDedent--
		return Token{Type: DEDENT, Literal: "", Line: l.line, Column: 0}
	}

	if l.IsAtLineStart() {
		if tok, cont := l.handleLineStart(); !cont {
			return tok
		} else if tok.Type != ILLEGAL {
			return tok
		}
	}

	l.skipWhitespace()

	tok := Token{Line: l.line, Column: l.column}

	if t, ok := l.handleSpecialChar(); ok {
		return t
	}

	switch l.ch {
	case '\n':
		tok = Token{Type: NEWLINE, Literal: string(l.ch), Line: l.line, Column: l.column}
		l.line++
		l.column = 0
		l.readChar()
		return tok
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
		return tok
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		return l.handleEOF()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
			return tok
		}
		tok = Token{Type: ILLEGAL, Literal: string(l.ch), Line: l.line, Column: l.column}
	}

	l.readChar()
	return tok
}

func (l *Lexer) handleLineStart() (Token, bool) {
	indent := 0
	for l.ch == ' ' {
		indent++
		l.readChar()
	}

	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.readChar()
		return Token{}, true
	}
	if l.ch == '#' {
		return Token{Type: COMMENT, Literal: l.readComment(), Line: l.line, Column: l.column}, false
	}

	currentIndent := l.indentStack[len(l.indentStack)-1]

	if indent > currentIndent {
		l.indentStack = append(l.indentStack, indent)
		return Token{Type: INDENT, Literal: "", Line: l.line, Column: 0}, false
	}
	if indent < currentIndent {
		for len(l.indentStack) > 1 && l.indentStack[len(l.indentStack)-1] > indent {
			l.indentStack = l.indentStack[:len(l.indentStack)-1]
			l.pendingDedent++
		}
		if l.pendingDedent > 0 {
			l.pendingDedent--
			return Token{Type: DEDENT, Literal: "", Line: l.line, Column: 0}, false
		}
	}

	return Token{Type: ILLEGAL}, true
}

var simpleTokens = map[byte]TokenType{
	':': COLON,
	',': COMMA,
	'=': EQUALS,
	'(': LPAREN,
	')': RPAREN,
	'$': DOLLAR,
}

// handleSpecialChar handles simple single-character tokens
func (l *Lexer) handleSpecialChar() (Token, bool) {
	if tokenType, ok := simpleTokens[l.ch]; ok {
		tok := Token{Type: tokenType, Literal: string(l.ch), Line: l.line, Column: l.column}
		l.readChar()
		return tok, true
	}
	return Token{}, false
}

// handleEOF handles end-of-file with pending dedents
func (l *Lexer) handleEOF() Token {
	for len(l.indentStack) > 1 {
		l.indentStack = l.indentStack[:len(l.indentStack)-1]
		l.pendingDedent++
	}
	if l.pendingDedent > 0 {
		l.pendingDedent--
		return Token{Type: DEDENT, Literal: "", Line: l.line, Column: 0}
	}
	return Token{Type: EOF, Literal: "", Line: l.line, Column: l.column}
}

func (l *Lexer) CheckIndentation() Token {
	return l.NextToken()
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	l.readChar()
	position := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	result := l.input[position:l.position]
	l.readChar()
	return result
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '-' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		if tok.Type == COMMENT {
			continue
		}
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return tokens
}

func (l *Lexer) IsAtLineStart() bool {
	return l.column == 1
}

func CountIndent(line string) int {
	count := 0
	for _, ch := range line {
		switch ch {
		case ' ':
			count++
		case '\t':
			count += 4
		default:
			return count
		}
	}
	return count
}

func StripIndent(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	minIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		indent := CountIndent(line)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return text
	}

	var result []string
	for _, line := range lines {
		if len(line) >= minIndent {
			result = append(result, line[minIndent:])
		} else {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}
