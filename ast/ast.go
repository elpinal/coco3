package ast

import "github.com/elpinal/coco3/token"

type Node interface {
	Pos() token.Pos // position of first character belonging to the node
	End() token.Pos // position of first character immediately after the node
}

type Expr interface {
	Node
	exprNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type (
	// A BadExpr node is a placeholder for expressions containing
	// syntax errors for which no correct expression nodes can be
	// created.
	//
	BadExpr struct {
		From, To token.Pos // position range of bad expression
	}

	// An Ident node represents an identifier.
	Ident struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name
	}

	BasicLit struct {
		ValuePos token.Pos   // literal position
		Kind     token.Token // token.STRING
		Value    string      // literal string; e.g. 'foo'
	}

	ParenExpr struct {
		Lparen token.Pos // position of "("
		Exprs  []Expr    // parenthesized expressions
		Rparen token.Pos // position of ")"
	}

	// A UnaryExpr node represents a unary expression.
	// Unary "*" expressions are represented via StarExpr nodes.
	//
	UnaryExpr struct {
		OpPos token.Pos   // position of Op
		Op    token.Token // operator
		X     Expr        // operand
	}
)

func (x *BadExpr) Pos() token.Pos   { return x.From }
func (x *Ident) Pos() token.Pos     { return x.NamePos }
func (x *BasicLit) Pos() token.Pos  { return x.ValuePos }
func (x *ParenExpr) Pos() token.Pos { return x.Lparen }
func (x *UnaryExpr) Pos() token.Pos { return x.OpPos }

func (x *BadExpr) End() token.Pos   { return x.To }
func (x *Ident) End() token.Pos     { return token.Pos(int(x.NamePos) + len(x.Name)) }
func (x *BasicLit) End() token.Pos  { return token.Pos(int(x.ValuePos) + len(x.Value)) }
func (x *ParenExpr) End() token.Pos { return x.Rparen + 1 }
func (x *UnaryExpr) End() token.Pos { return x.X.End() }

func (*BadExpr) exprNode()   {}
func (*Ident) exprNode()     {}
func (*BasicLit) exprNode()  {}
func (*ParenExpr) exprNode() {}
func (*UnaryExpr) exprNode() {}

func (id *Ident) String() string {
	if id != nil {
		return id.Name
	}
	return "<nil>"
}

type (
	// A BadStmt node is a placeholder for statements containing
	// syntax errors for which no correct statement nodes can be
	// created.
	//
	BadStmt struct {
		From, To token.Pos // position range of bad statement
	}

	ExecStmt struct {
		Args []Expr
	}
)

func (s *BadStmt) Pos() token.Pos  { return s.From }
func (s *ExecStmt) Pos() token.Pos { return s.Args[0].Pos() }

func (s *BadStmt) End() token.Pos  { return s.To }
func (s *ExecStmt) End() token.Pos { return s.Args[len(s.Args)-1].End() }

func (*BadStmt) stmtNode()  {}
func (*ExecStmt) stmtNode() {}

type File struct {
	Name  *Ident // package name
	Lines []Stmt
}

func (f *File) Pos() token.Pos { return token.Pos(1) }
func (f *File) End() token.Pos {
	if n := len(f.Lines); n > 0 {
		return f.Lines[n-1].End()
	}
	return token.Pos(1)
}
