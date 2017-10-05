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
%type <list> string_list

%token <token> ILLEGAL

%token <token> IDENT STRING LBRACK RBRACK NUM COLON

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
        | string_list
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

string_list:
        LBRACK RBRACK
        {
                $$ = &ast.Empty{}
        }
        | STRING COLON string_list
        {
                $$ = &ast.Cons{Head: $1.Lit, Tail: $3}
        }

%%
