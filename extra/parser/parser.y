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
        args    []token.Token
}

%type <command> top command
%type <args> args

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
        IDENT args
        {
                $$ = &ast.Command{$1.Lit, $2}
        }

args:
        {
                $$ = []token.Token{}
        }
        | args STRING
        {
                $$ = append($1, $2)
        }

%%
