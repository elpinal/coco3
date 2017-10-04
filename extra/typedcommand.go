package extra

type TypedCommand struct {
	params []Type
	fn     func([]string) error
}

type Type int

const (
	String Type = iota + 1
)
