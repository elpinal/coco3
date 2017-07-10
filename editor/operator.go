package editor

const (
	OpNop = iota
	OpDelete
	OpYank
	OpChange
	OpLower
	OpUpper
)

var opChars = map[string]int{
	"d":  OpDelete,
	"y":  OpYank,
	"c":  OpChange,
	"gu": OpLower,
	"gU": OpUpper,
}
