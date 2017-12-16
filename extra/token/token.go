package token

type Token struct {
	Lit string

	Line   uint
	Column uint
}
