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
	PROMPT
	NOTIFY

	SHELL
	DOLLAR

	TEMPLATE
	GROUP
	BEFORE
	AFTER
	ALIAS
	EXTENDS
)

var keywords = map[string]TokenType{
	"var":          VAR,
	"task":         TASK,
	"profile":      PROFILE,
	"include":      INCLUDE,
	"deps":         DEPS,
	"run":          RUN,
	"env":          ENV,
	"watch":        WATCH,
	"matrix":       MATRIX,
	"docker":       DOCKER,
	"remote":       REMOTE,
	"desc":         DESC,
	"parallel":     PARALLEL,
	"if":           IF,
	"cache":        CACHE,
	"inputs":       INPUTS,
	"outputs":      OUTPUTS,
	"ignore":       IGNORE,
	"profile_task": PROFILE_TASK,
	"secrets":      SECRETS,
	"pre":          PRE,
	"retries":      RETRIES,
	"retry_delay":  RETRY_DELAY,
	"timeout":      TIMEOUT,
	"prompt":       PROMPT,
	"notify":       NOTIFY,
	"shell":        SHELL,
	"template":     TEMPLATE,
	"group":        GROUP,
	"before":       BEFORE,
	"after":        AFTER,
	"alias":        ALIAS,
	"extends":      EXTENDS,
	"true":         IDENT,
	"false":        IDENT,
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

var tokenNames = map[TokenType]string{
	ILLEGAL:      "ILLEGAL",
	EOF:          "EOF",
	COMMENT:      "COMMENT",
	IDENT:        "IDENT",
	STRING:       "STRING",
	NUMBER:       "NUMBER",
	NEWLINE:      "NEWLINE",
	INDENT:       "INDENT",
	DEDENT:       "DEDENT",
	VAR:          "VAR",
	TASK:         "TASK",
	PROFILE:      "PROFILE",
	INCLUDE:      "INCLUDE",
	COLON:        "COLON",
	COMMA:        "COMMA",
	EQUALS:       "EQUALS",
	LPAREN:       "LPAREN",
	RPAREN:       "RPAREN",
	DEPS:         "DEPS",
	RUN:          "RUN",
	ENV:          "ENV",
	WATCH:        "WATCH",
	MATRIX:       "MATRIX",
	DOCKER:       "DOCKER",
	REMOTE:       "REMOTE",
	DESC:         "DESC",
	PARALLEL:     "PARALLEL",
	IF:           "IF",
	CACHE:        "CACHE",
	INPUTS:       "INPUTS",
	OUTPUTS:      "OUTPUTS",
	IGNORE:       "IGNORE",
	PROFILE_TASK: "PROFILE_TASK",
	SECRETS:      "SECRETS",
	PRE:          "PRE",
	RETRIES:      "RETRIES",
	RETRY_DELAY:  "RETRY_DELAY",
	TIMEOUT:      "TIMEOUT",
	PROMPT:       "PROMPT",
	NOTIFY:       "NOTIFY",
	SHELL:        "SHELL",
	DOLLAR:       "DOLLAR",
	TEMPLATE:     "TEMPLATE",
	GROUP:        "GROUP",
	BEFORE:       "BEFORE",
	AFTER:        "AFTER",
	ALIAS:        "ALIAS",
	EXTENDS:      "EXTENDS",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}
