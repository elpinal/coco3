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
        exprs   []ast.Expr
        expr    ast.Expr
        list    ast.List
}

%type <command> top command
%type <exprs> exprs
%type <expr> expr
%type <list> empty

%token <token> ILLEGAL

%token <token> IDENT STRING LBRACK RBRACK NUM

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
        IDENT exprs
        {
                $$ = &ast.Command{$1.Lit, $2}
        }

expr:
        STRING
        {
                $$ = &ast.String{$1.Lit}
        }
        | NUM
        {
                $$ = &ast.Int{$1.Lit}
        }
        | empty
        {
                $$ = $1
        }

exprs:
        {
                $$ = nil
        }
        | exprs expr
        {
                $$ = append($1, $2)
        }

empty:
        LBRACK RBRACK
        {
                $$ = &ast.Empty{}
        }

%%
