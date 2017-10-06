package parser

//go:generate goyacc -o parser.go parser.y

import (
	"bytes"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/token"
	"github.com/elpinal/coco3/extra/types"
)

const eof = 0

type exprLexer struct {
	src []byte // source
	ch  rune   // current character
	err *ParseError

	// result
	expr *ast.Command

	// information for error messages
	off    uint // start at 0
	line   uint // start at 1
	column uint // start at 1

	// information for current token
	tokLine   uint
	tokColumn uint
}

func newLexer(src []byte) *exprLexer {
	l := &exprLexer{
		src:  src,
		line: 1,
	}
	l.next()
	return l
}

func isAlphabet(c rune) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

func isIdent(c rune) bool {
	return isAlphabet(c) || isNumber(c) || c == '-'
}

func isNumber(c rune) bool {
	return '0' <= c && c <= '9'
}

func isQuote(c rune) bool {
	return c == '\''
}

func (x *exprLexer) Lex(yylval *yySymType) int {
	for {
		x.tokLine = x.line
		x.tokColumn = x.column
		c := x.ch
		switch c {
		case eof:
			return eof
		case ' ':
			x.next()
		case '\'':
			x.next()
			return x.str(yylval)
		case '[':
			x.next()
			return LBRACK
		case ']':
			x.next()
			return RBRACK
		case ':':
			x.next()
			return COLON
		case ',':
			x.next()
			return COMMA
		default:
			if isAlphabet(c) {
				return x.ident(yylval)
			}
			if isNumber(c) {
				return x.num(yylval)
			}
			fmt.Fprintf(os.Stderr, "%d:%d: invalid character: %[3]U %[3]q\n", x.line, x.column, c)
			return ILLEGAL
		}
	}
}

func (x *exprLexer) ident(yylval *yySymType) int {
	x.takeWhile(types.Ident, isIdent, yylval)
	return IDENT
}

func (x *exprLexer) str(yylval *yySymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			x.err = &ParseError{
				Line:   x.line,
				Column: x.column,
				Msg:    fmt.Sprintf("WriteRune: %s", err),
			}
		}
	}
	var b bytes.Buffer
	for !isQuote(x.ch) {
		if x.ch == eof {
			x.err = &ParseError{
				Line:   x.line,
				Column: x.column,
				Msg:    "string literal not terminated: unexpected EOF",
			}
			return STRING
		}
		add(&b, x.ch)
		x.next()
	}
	yylval.token = token.Token{
		Kind: types.String,
		Lit:  b.String(),
	}
	x.next()
	return STRING
}

func (x *exprLexer) num(yylval *yySymType) int {
	x.takeWhile(types.Int, isNumber, yylval)
	return NUM
}

func (x *exprLexer) takeWhile(kind types.Type, f func(rune) bool, yylval *yySymType) {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			x.err = &ParseError{
				Line:   x.line,
				Column: x.column,
				Msg:    fmt.Sprintf("WriteRune: %s", err),
			}
		}
	}
	var b bytes.Buffer
	for f(x.ch) && x.ch != eof {
		add(&b, x.ch)
		x.next()
	}
	yylval.token = token.Token{
		Kind: kind,
		Lit:  b.String(),
	}
}

func (x *exprLexer) next() {
	if len(x.src) == 0 {
		x.ch = eof
		return
	}
	c, size := utf8.DecodeRune(x.src)
	x.src = x.src[size:]
	x.off++
	if c == '\n' {
		x.line++
		x.column = 0
	} else {
		x.column++
	}
	if c == utf8.RuneError && size == 1 {
		x.err = &ParseError{
			Line:   x.line,
			Column: x.column,
			Msg:    "next: invalid utf8",
		}
		x.next()
		return
	}
	x.ch = c
}

func (x *exprLexer) Error(s string) {
	x.err = &ParseError{
		Line:   x.tokLine,
		Column: x.tokColumn,
		Msg:    s,
	}
}

func Parse(src []byte) (*ast.Command, error) {
	l := newLexer(src)
	yyErrorVerbose = true
	yyParse(l)
	if l.err != nil {
		l.err.Src = string(src)
		return nil, l.err
	}
	return l.expr, nil
}
