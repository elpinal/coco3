package ast

import "github.com/elpinal/coco3/extra/token"

type Command struct {
	Name string
	Args []token.Token
}

type List interface {
	Length() int
}

type Empty struct{}

func (e *Empty) Length() int {
	return 0
}

type Cons struct {
	Head string
	Tail List
}

func (c *Cons) Length() int {
	return 1 + c.Tail.Length()
}
