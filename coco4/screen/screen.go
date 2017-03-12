package screen

type Screen interface {
	Refresh(string, []rune, int)
	SetLastLine(string)
}
