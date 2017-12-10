package parser

//go:generate goyacc -o parser.go parser.y

import (
	"bytes"
	"fmt"
	"unicode/utf8"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/token"
	"github.com/elpinal/coco3/extra/types"
)

const eof = 0

type exprLexer struct {
	src []byte // source
	r   rune   // current character

	off    uint // starts from 0
	line   uint // starts from 1
	column uint // starts from 1

	// information for the start position of current token
	tokLine   uint
	tokColumn uint

	// result
	expr *ast.Command

	// channel for error
	errCh chan *ParseError
}

func newLexer(src []byte) *exprLexer {
	l := &exprLexer{
		src:   src,
		line:  1,
		errCh: make(chan *ParseError),
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

func (l *exprLexer) emitError(format string, args ...interface{}) {
	select {
	case l.errCh <- l.errorAtHere(format, args...):
	default:
		// If errCh is blocked (i.e. another error had occurred), this
		// error message is ignored.
	}
}

func (l *exprLexer) errorAtHere(format string, args ...interface{}) *ParseError {
	return &ParseError{
		Line:   l.line,
		Column: l.column,
		Msg:    fmt.Sprintf(format, args...),
	}
}

func (l *exprLexer) Lex(yylval *yySymType) int {
	for {
		l.tokLine = l.line
		l.tokColumn = l.column
		c := l.r
		// TODO: set yylval.token on all cases.
		switch c {
		case eof:
			return eof
		case ' ', '\n':
			l.next()
		case '\'':
			l.next()
			return l.str(yylval)
		case '[':
			l.next()
			return LBRACK
		case ']':
			l.next()
			return RBRACK
		case ':':
			l.next()
			return COLON
		case ',':
			l.next()
			return COMMA
		case '!':
			l.next()
			yylval.token = token.Token{
				Kind:   types.Ident,
				Lit:    "!",
				Line:   l.tokLine,
				Column: l.tokColumn,
			}
			return int(c)
		default:
			if isAlphabet(c) {
				return l.ident(yylval)
			}
			if isNumber(c) {
				return l.num(yylval)
			}
			l.emitError("invalid character: %[1]U %[1]q", c)
			return ILLEGAL
		}
	}
}

func (l *exprLexer) ident(yylval *yySymType) int {
	l.takeWhile(types.Ident, isIdent, yylval)
	return IDENT
}

func (l *exprLexer) str(yylval *yySymType) int {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			l.emitError("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	for !isQuote(l.r) {
		if l.r == eof {
			l.errCh <- &ParseError{
				Line:   l.tokLine,
				Column: l.tokColumn,
				Msg:    "string literal not terminated: unexpected EOF",
			}
			return STRING
		}
		if l.r == '\\' {
			line := l.line
			column := l.column
			l.next()
			switch l.r {
			case '\'', '\\':
				add(&b, l.r)
				l.next()
			case 'n':
				add(&b, '\n')
				l.next()
			case eof:
				l.errCh <- &ParseError{
					Line:   line,
					Column: column,
					Msg:    "string literal not terminated: unexpected EOF",
				}
				return STRING
			default:
				l.errCh <- &ParseError{
					Line:   line,
					Column: column,
					Msg:    fmt.Sprintf("unknown escape sequence: \\%c", l.r),
				}
				l.next()
				return STRING
			}
			continue
		}
		add(&b, l.r)
		l.next()
	}
	yylval.token = token.Token{
		Kind:   types.String,
		Lit:    b.String(),
		Line:   l.tokLine,
		Column: l.tokColumn,
	}
	l.next()
	return STRING
}

func (l *exprLexer) num(yylval *yySymType) int {
	l.takeWhile(types.Int, isNumber, yylval)
	return NUM
}

func (l *exprLexer) takeWhile(kind types.Type, f func(rune) bool, yylval *yySymType) {
	add := func(b *bytes.Buffer, c rune) {
		if _, err := b.WriteRune(c); err != nil {
			l.emitError("WriteRune: %s", err)
		}
	}
	var b bytes.Buffer
	for f(l.r) && l.r != eof {
		add(&b, l.r)
		l.next()
	}
	yylval.token = token.Token{
		Kind:   kind,
		Lit:    b.String(),
		Line:   l.tokLine,
		Column: l.tokColumn,
	}
}

func (l *exprLexer) next() {
	if len(l.src) == 0 {
		l.r = eof
		return
	}
	c, size := utf8.DecodeRune(l.src)
	l.src = l.src[size:]
	l.off++
	if c == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	if c == utf8.RuneError && size == 1 {
		l.emitError("next: invalid utf8")
		l.next()
		return
	}
	l.r = c
}

func (l *exprLexer) Error(s string) {
	l.errCh <- &ParseError{
		Line:   l.tokLine,
		Column: l.tokColumn,
		Msg:    s,
	}
}

func (l *exprLexer) run() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		yyParse(l)
		done <- struct{}{}
	}()
	return done
}

func init() {
	yyErrorVerbose = true
}

func Parse(src []byte) (*ast.Command, error) {
	l := newLexer(src)
	done := l.run()
	select {
	case err := <-l.errCh:
		err.Src = string(src)
		return nil, err
	case <-done:
	}
	return l.expr, nil
}
