package extra

type TypedCommand struct {
	params []Type
	fn     func(string)
}

type Type int

const (
	String Type = iota + 1
)
