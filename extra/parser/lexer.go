package parser

//go:generate goyacc -o parser.go parser.y

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/token"
	"github.com/elpinal/coco3/extra/typed"
)

const eof = 0

type exprLexer struct {
	src []byte // source
	ch  rune   // current character
	err error

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

func isNotQuote(c rune) bool {
	return c != '\''
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
		default:
			if isAlphabet(c) {
				return x.ident(yylval)
			}
			fmt.Fprintf(os.Stderr, "[%d:%d]: invalid character: %[3]U %[3]q\n", x.line, x.column, c)
			return ILLEGAL
		}
	}
}

func (x *exprLexer) ident(yylval *yySymType) int {
	x.takeWhile(typed.Ident, isAlphabet, yylval)
	return IDENT
}

func (x *exprLexer) str(yylval *yySymType) int {
	x.takeWhile(typed.String, isNotQuote, yylval)
	x.next()
	return STRING
}

func (x *exprLexer) takeWhile(kind typed.Type, f func(rune) bool, yylval *yySymType) {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			x.err = fmt.Errorf("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	for f(x.ch) {
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
		x.err = errors.New("next: invalid utf8")
		x.next()
		return
	}
	x.ch = c
}

func (x *exprLexer) Error(s string) {
	x.err = errors.New(s)
}

func Parse(src []byte) (*ast.Command, error) {
	l := newLexer(src)
	yyErrorVerbose = true
	yyParse(l)
	if l.err != nil {
		return nil, l.err
	}
	return l.expr, nil
}
