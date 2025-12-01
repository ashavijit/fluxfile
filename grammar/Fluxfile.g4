grammar Fluxfile;

fluxfile
    : statement* EOF
    ;

statement
    : varDecl
    | taskDecl
    | profileDecl
    | includeDecl
    | NEWLINE
    ;

varDecl
    : VAR IDENT EQUALS expr NEWLINE
    ;

taskDecl
    : TASK IDENT COLON NEWLINE INDENT taskBody DEDENT
    ;

taskBody
    : taskDirective+
    ;

taskDirective
    : descDirective
    | depsDirective
    | parallelDirective
    | ifDirective
    | runDirective
    | envDirective
    | watchDirective
    | ignoreDirective
    | matrixDirective
    | cacheDirective
    | inputsDirective
    | outputsDirective
    | dockerDirective
    | remoteDirective
    ;

descDirective
    : DESC COLON (STRING | identList) NEWLINE
    ;

depsDirective
    : DEPS COLON identList NEWLINE
    ;

parallelDirective
    : PARALLEL COLON IDENT NEWLINE
    ;

ifDirective
    : IF COLON conditionExpr NEWLINE
    ;

runDirective
    : RUN COLON NEWLINE INDENT commandList DEDENT
    ;

envDirective
    : ENV COLON NEWLINE INDENT envPairList DEDENT
    ;

watchDirective
    : WATCH COLON pattern NEWLINE
    ;

ignoreDirective
    : IGNORE COLON NEWLINE INDENT patternList DEDENT
    ;

matrixDirective
    : MATRIX COLON NEWLINE INDENT matrixDimensions DEDENT
    ;

cacheDirective
    : CACHE COLON IDENT NEWLINE
    ;

inputsDirective
    : INPUTS COLON NEWLINE INDENT patternList DEDENT
    ;

outputsDirective
    : OUTPUTS COLON NEWLINE INDENT patternList DEDENT
    ;

dockerDirective
    : DOCKER COLON IDENT NEWLINE
    ;

remoteDirective
    : REMOTE COLON STRING NEWLINE
    ;

profileDecl
    : PROFILE IDENT COLON NEWLINE INDENT envDirective DEDENT
    ;

includeDecl
    : INCLUDE STRING NEWLINE
    ;

commandList
    : command+
    ;

command
    : ~(NEWLINE | DEDENT)+ NEWLINE
    ;

envPairList
    : envPair+
    ;

envPair
    : IDENT EQUALS expr NEWLINE
    ;

matrixDimensions
    : matrixDimension+
    ;

matrixDimension
    : IDENT COLON identList NEWLINE
    ;

patternList
    : pattern+
    ;

pattern
    : (IDENT | STRING) NEWLINE
    ;

identList
    : IDENT (COMMA IDENT)*
    ;

conditionExpr
    : IDENT compOp expr
    ;

compOp
    : '=='
    | '!='
    | '>'
    | '<'
    | '>='
    | '<='
    ;

expr
    : STRING
    | NUMBER
    | IDENT
    | shellExpr
    ;

shellExpr
    : DOLLAR LPAREN SHELL STRING RPAREN
    ;

VAR         : 'var' ;
TASK        : 'task' ;
PROFILE     : 'profile' ;
INCLUDE     : 'include' ;
DESC        : 'desc' ;
DEPS        : 'deps' ;
PARALLEL    : 'parallel' ;
IF          : 'if' ;
RUN         : 'run' ;
ENV         : 'env' ;
WATCH       : 'watch' ;
IGNORE      : 'ignore' ;
MATRIX      : 'matrix' ;
CACHE       : 'cache' ;
INPUTS      : 'inputs' ;
OUTPUTS     : 'outputs' ;
DOCKER      : 'docker' ;
REMOTE      : 'remote' ;
SHELL       : 'shell' ;

COLON       : ':' ;
COMMA       : ',' ;
EQUALS      : '=' ;
LPAREN      : '(' ;
RPAREN      : ')' ;
DOLLAR      : '$' ;

INDENT      : '<INDENT>' ;
DEDENT      : '<DEDENT>' ;
NEWLINE     : '\r'? '\n' ;

STRING      : '"' (~["\r\n])* '"' ;
NUMBER      : [0-9]+ ('.' [0-9]+)? ;
IDENT       : [a-zA-Z_][a-zA-Z0-9_-]* ;

WS          : [ \t]+ -> skip ;
COMMENT     : '#' ~[\r\n]* -> skip ;
