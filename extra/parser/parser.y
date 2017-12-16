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
        def     ast.Def
}

%type <command> top command
%type <exprs> exprs
%type <expr> expr
%type <list> string_list sep_by_commas
%type <def> def

%token <token> ILLEGAL

%token <token> IDENT STRING LBRACK RBRACK NUM COLON COMMA '!'
%token <token> DEF
%token <token> '='
%token <token> FN

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
        {
                $$ = nil
        }
        |
        IDENT exprs
        {
                $$ = &ast.Command{$1, $2}
        }
        | '!' IDENT exprs
        {
                $$ = &ast.Command{
                        token.Token{Lit: "exec", Line: $1.Line, Column: $1.Column},
                        append([]ast.Expr{&ast.String{$2.Lit}}, $3...),
                }
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
        | IDENT
        {
                $$ = &ast.Ident{$1.Lit}
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
        | LBRACK sep_by_commas RBRACK
        {
                $$ = $2
        }

sep_by_commas:
        STRING
        {
                $$ = &ast.Cons{Head: $1.Lit, Tail: &ast.Empty{}}
        }
        | STRING COMMA sep_by_commas
        {
                $$ = &ast.Cons{Head: $1.Lit, Tail: $3}
        }

def:
        DEF IDENT '=' expr
        {
                $$ = &ast.Def{Name: $2, Expr: $4}
        }

%%
