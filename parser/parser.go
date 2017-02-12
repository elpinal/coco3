package parser

import (
	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/scanner"
	"github.com/elpinal/coco3/token"
)

// The parser structure holds the parser's internal state.
type parser struct {
	file    *token.File
	errors  scanner.ErrorList
	scanner scanner.Scanner

	// Next token
	pos token.Pos   // token position
	tok token.Token // one token look-ahead
	lit string      // token literal

}

func (p *parser) init(fset *token.FileSet, filename string, src []byte) {
	p.file = fset.AddFile(filename, -1, len(src))
	eh := func(pos token.Position, msg string) { p.errors.Add(pos, msg) }
	p.scanner.Init(p.file, src, eh)

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
	epos := p.file.Position(pos)
	p.errors.Add(epos, msg)
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

func (p *parser) checkExpr(x ast.Expr) ast.Expr {
	switch x.(type) {
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.ParenExpr:
		panic("unreachable")
	case *ast.UnaryExpr:
	default:
		// all other nodes are not proper expressions
		p.errorExpected(x.Pos(), "expression")
		x = &ast.BadExpr{From: x.Pos(), To: p.safePos(x.End())}
	}
	return x
}

func (p *parser) safePos(pos token.Pos) (res token.Pos) {
	defer func() {
		if recover() != nil {
			res = token.Pos(p.file.Base() + p.file.Size()) // EOF position
		}
	}()
	_ = p.file.Offset(pos) // trigger a panic if position is out-of-range
	return pos
}

func (p *parser) parseUnary() ast.Expr {
	pos, op := p.pos, p.tok
	p.next()
	x := p.parseExpr()
	return &ast.UnaryExpr{OpPos: pos, Op: op, X: p.checkExpr(x)}
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

func (p *parser) parseList() *ast.ParenExpr {
	pos := p.pos
	var list []ast.Expr
	p.next()
	for p.tok != token.RPAREN && p.tok != token.EOF {
		expr := p.parseExpr()
		list = append(list, expr)
	}
	p.next()
	return &ast.ParenExpr{Lparen: pos, Exprs: list, Rparen: p.pos - 1}
}

func (p *parser) parseExpr() ast.Expr {
	switch p.tok {
	case token.LPAREN:
		return p.parseList()
	case token.IDENT:
		return p.parseIdent()
	case token.REDIRIN, token.REDIROUT:
		return p.parseUnary()
	}
	pos := p.pos
	p.errorExpected(pos, "expression")
	p.next()
	return &ast.BadExpr{From: pos, To: p.pos}
}

func (p *parser) parseLine() ast.Stmt {
	var args []ast.Expr
	for p.tok != token.SEMICOLON && p.tok != token.EOF {
		args = append(args, p.parseExpr())
	}
	p.next()
	return &ast.ExecStmt{Args: args}
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
