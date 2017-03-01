package editor

const (
	OpNop = iota
	OpDelete
	OpYank
)

var opChars = map[rune]int{
	'd': OpDelete,
	'y': OpYank,
}
