package ast

import "github.com/elpinal/coco3/extra/token"

type Command struct {
	Name string
	Args  []token.Token
}
