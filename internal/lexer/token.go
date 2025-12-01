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
	DESC
	PARALLEL
	IF
	CACHE
	INPUTS
	OUTPUTS
	IGNORE
	PROFILE_TASK
	SECRETS
	PRE
	RETRIES
	RETRY_DELAY
	TIMEOUT

	SHELL
	DOLLAR
)

var keywords = map[string]TokenType{
	"var":         VAR,
	"task":        TASK,
	"profile":     PROFILE,
	"include":     INCLUDE,
	"deps":        DEPS,
	"run":         RUN,
	"env":         ENV,
	"watch":       WATCH,
	"matrix":      MATRIX,
	"docker":      DOCKER,
	"remote":      REMOTE,
	"desc":        DESC,
	"parallel":    PARALLEL,
	"if":          IF,
	"cache":       CACHE,
	"inputs":      INPUTS,
	"outputs":     OUTPUTS,
	"ignore":      IGNORE,
	"profile_task": PROFILE_TASK,
	"secrets":     SECRETS,
	"pre":         PRE,
	"retries":     RETRIES,
	"retry_delay": RETRY_DELAY,
	"timeout":     TIMEOUT,
	"shell":       SHELL,
	"true":        IDENT,
	"false":       IDENT,
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
	case DESC:
		return "DESC"
	case PARALLEL:
		return "PARALLEL"
	case IF:
		return "IF"
	case CACHE:
		return "CACHE"
	case INPUTS:
		return "INPUTS"
	case OUTPUTS:
		return "OUTPUTS"
	case IGNORE:
		return "IGNORE"
	case PROFILE_TASK:
		return "PROFILE_TASK"
	case SECRETS:
		return "SECRETS"
	case PRE:
		return "PRE"
	case RETRIES:
		return "RETRIES"
	case RETRY_DELAY:
		return "RETRY_DELAY"
	case TIMEOUT:
		return "TIMEOUT"
	case SHELL:
		return "SHELL"
	case DOLLAR:
		return "DOLLAR"
	default:
		return "UNKNOWN"
	}
}
