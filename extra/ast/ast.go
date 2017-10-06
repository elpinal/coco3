package ast

import (
	"fmt"

	"github.com/elpinal/coco3/extra/token"
	"github.com/elpinal/coco3/extra/types"
)

type Command struct {
	Name token.Token
	Args []Expr
}

// Expressions

type Expr interface {
	Expr()
	Type() types.Type
}

func (_ *String) Expr() {}
func (_ *Int) Expr()    {}
func (_ *Ident) Expr()  {}
func (_ *Empty) Expr()  {}
func (_ *Cons) Expr()   {}

// Simple types

type (
	String struct {
		Lit string
	}

	Int struct {
		Lit string
	}

	Ident struct {
		Lit string
	}
)

func (_ *String) Type() types.Type {
	return types.String
}

func (_ *Int) Type() types.Type {
	return types.Int
}

func (_ *Ident) Type() types.Type {
	return types.Ident
}

func (s *String) String() string {
	return fmt.Sprintf("%q", s.Lit)
}

func (i *Int) String() string {
	return i.Lit
}

func (id *Ident) String() string {
	return fmt.Sprintf("%q", id.Lit)
}

// Lists

type List interface {
	Length() int
	Expr
}

type (
	Cons struct {
		Head string
		Tail List
	}

	Empty struct{}
)

func (e *Empty) Type() types.Type {
	return types.StringList
}

func (c *Cons) Type() types.Type {
	return types.StringList
}

func (e *Empty) Length() int {
	return 0
}

func (c *Cons) Length() int {
	return 1 + c.Tail.Length()
}
