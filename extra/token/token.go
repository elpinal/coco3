package token

import "github.com/elpinal/coco3/extra/types"

type Token struct {
	Kind
	Lit string

	Line   uint
	Column uint
}

type Kind struct {
	types.Type
}

func KindOf(t types.Type) Kind {
	return Kind{Type: t}
}
