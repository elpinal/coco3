package typed

type Command struct {
	Params []Type
	Fn     func([]string) error
}

type Type int

const (
	String Type = iota + 1
)
