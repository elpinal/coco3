package extra

type TypedCommand struct {
	params []Type
}

type Type int

const (
	String Type = iota + 1
)
