package ast

import "github.com/elpinal/coco3/extra/types"

type Command struct {
	Name string
	Args []Expr
}

type Expr interface {
	Expr()
	Type() types.Type
}

func (_ *String) Expr() {}
func (_ *Empty) Expr()  {}
func (_ *Cons) Expr()   {}

type String struct {
	Lit string
}

func (s *String) Type() types.Type {
	return types.String
}

type List interface {
	Length() int
	Expr
}

type Empty struct{}

func (e *Empty) Type() types.Type {
	return types.StringList
}

func (e *Empty) Length() int {
	return 0
}

type Cons struct {
	Head string
	Tail List
}

func (c *Cons) Type() types.Type {
	return types.StringList
}

func (c *Cons) Length() int {
	return 1 + c.Tail.Length()
}
