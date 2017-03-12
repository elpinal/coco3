package editor

const (
	OpNop = iota
	OpDelete
	OpYank
	OpChange
)

var opChars = map[rune]int{
	'd': OpDelete,
	'y': OpYank,
	'c': OpChange,
}
