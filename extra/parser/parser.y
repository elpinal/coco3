%{

package parser

import "github.com/elpinal/coco3/extra/token"

%}

%union {
        token token.Token
}

%type <token> top

%token <token> ILLEGAL

%token <token> IDENT

%%

top:
        IDENT
        {
                $$ = $1
                if l, ok := yylex.(*exprLexer); ok {
                        l.expr = $$
                }
        }

%%
