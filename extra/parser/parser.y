%{

package parser

import (
        "github.com/elpinal/coco3/extra/ast"
        "github.com/elpinal/coco3/extra/token"
)

%}

%union {
        token   token.Token
        command *ast.Command
}

%type <command> top command

%token <token> ILLEGAL

%token <token> IDENT STRING

%%

top:
        command
        {
                $$ = $1
                if l, ok := yylex.(*exprLexer); ok {
                        l.expr = $$
                }
        }

command:
        IDENT STRING
        {
                $$ = &ast.Command{$1.Lit, $2}
        }

%%
