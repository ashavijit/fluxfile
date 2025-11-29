package lexer

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	IDENT
	STRING
	NUMBER

	NEWLINE
	INDENT
	DEDENT

	VAR
	TASK
	PROFILE
	INCLUDE

	COLON
	COMMA
	EQUALS
	LPAREN
	RPAREN

	DEPS
	RUN
	ENV
	WATCH
	MATRIX
	DOCKER
	REMOTE

	SHELL
	DOLLAR
)

var keywords = map[string]TokenType{
	"var":     VAR,
	"task":    TASK,
	"profile": PROFILE,
	"include": INCLUDE,
	"deps":    DEPS,
	"run":     RUN,
	"env":     ENV,
	"watch":   WATCH,
	"matrix":  MATRIX,
	"docker":  DOCKER,
	"remote":  REMOTE,
	"shell":   SHELL,
	"true":    IDENT,
	"false":   IDENT,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case COMMENT:
		return "COMMENT"
	case IDENT:
		return "IDENT"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case NEWLINE:
		return "NEWLINE"
	case INDENT:
		return "INDENT"
	case DEDENT:
		return "DEDENT"
	case VAR:
		return "VAR"
	case TASK:
		return "TASK"
	case PROFILE:
		return "PROFILE"
	case INCLUDE:
		return "INCLUDE"
	case COLON:
		return "COLON"
	case COMMA:
		return "COMMA"
	case EQUALS:
		return "EQUALS"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case DEPS:
		return "DEPS"
	case RUN:
		return "RUN"
	case ENV:
		return "ENV"
	case WATCH:
		return "WATCH"
	case MATRIX:
		return "MATRIX"
	case DOCKER:
		return "DOCKER"
	case REMOTE:
		return "REMOTE"
	case SHELL:
		return "SHELL"
	case DOLLAR:
		return "DOLLAR"
	default:
		return "UNKNOWN"
	}
}
