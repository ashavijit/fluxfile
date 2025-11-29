package lexer

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	input := `var PROJECT = flux

task build:
    run:
        go build -o app
`

	l := New(input)
	tokens := l.Tokenize()

	if len(tokens) == 0 {
		t.Fatal("Expected tokens, got none")
	}

	if tokens[0].Type != VAR {
		t.Errorf("Expected VAR token, got %s", tokens[0].Type)
	}
}

func TestShellExpr(t *testing.T) {
	input := `var VERSION = $(shell "git describe --tags")`

	l := New(input)
	tokens := l.Tokenize()

	found := false
	for _, tok := range tokens {
		if tok.Type == SHELL {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected to find SHELL token")
	}
}

func TestIndentation(t *testing.T) {
	input := `task test:
    run:
        echo hello
`

	l := New(input)
	tokens := l.Tokenize()

	indentCount := 0
	for _, tok := range tokens {
		if tok.Type == INDENT {
			indentCount++
		}
	}

	if indentCount == 0 {
		t.Error("Expected INDENT tokens")
	}
}
