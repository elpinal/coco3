package token

type Pos int

const NoPos Pos = 0

func (p Pos) IsValid() bool {
	return p != NoPos
}
