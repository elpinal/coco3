package parser

import (
	"log"

	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/scanner"
	"github.com/elpinal/coco3/token"
)

// The parser structure holds the parser's internal state.
type parser struct {
	scanner scanner.Scanner

	// Next token
	pos token.Pos   // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

}

func (p *parser) init(src []byte) {
	p.scanner.Init(src)

	p.next()
}

// ----------------------------------------------------------------------------
// Parsing support

// Advance to the next token.
func (p *parser) next0() {
	p.pos, p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) next() {
	p.next0()
}

func (p *parser) error(pos token.Pos, msg string) {
	log.Println("parser.error:", msg)
}

func (p *parser) errorExpected(pos token.Pos, msg string) {
	msg = "expected " + msg
	if pos == p.pos {
		// the error happened at the current position;
		// make the error message more specific
		if p.tok == token.SEMICOLON && p.lit == "\n" {
			msg += ", found newline"
		} else {
			msg += ", found '" + p.tok.String() + "'"
			if p.tok.IsLiteral() {
				msg += " " + p.lit
			}
		}
	}
	p.error(pos, msg)
}

func (p *parser) expect(tok token.Token) token.Pos {
	pos := p.pos
	if p.tok != tok {
		p.errorExpected(pos, "'"+tok.String()+"'")
	}
	p.next() // make progress
	return pos
}

func (p *parser) parseIdent() *ast.Ident {
	pos := p.pos
	name := "_"
	if p.tok == token.IDENT {
		name = p.lit
		p.next()
	} else {
		p.expect(token.IDENT) // use expect() error handling
	}
	return &ast.Ident{NamePos: pos, Name: name}
}

func (p *parser) parseLine() ast.Stmt {
	cmd := p.parseIdent()
	var args []ast.Expr
	for p.tok != token.SEMICOLON && p.tok != token.EOF {
		args = append(args, p.parseIdent())
	}
	p.next()
	return &ast.ExecStmt{Cmd: cmd, Args: args}
}

// ----------------------------------------------------------------------------
// Source files

func (p *parser) parseFile() *ast.File {
	var lines []ast.Stmt
	for p.tok != token.EOF {
		lines = append(lines, p.parseLine())
	}

	return &ast.File{
		Lines: lines,
	}
}
