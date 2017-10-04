package token

import "github.com/elpinal/coco3/extra/typed"

type Token struct {
	Kind typed.Type
	Lit  string
}
