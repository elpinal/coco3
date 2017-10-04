%{

package parser

%}

%union {
        token int
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
