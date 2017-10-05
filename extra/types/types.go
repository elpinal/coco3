package types

type Type int

const (
	String Type = iota + 1
	Int
	Ident
	StringList
)

func (t Type) String() string {
	switch t {
	case String:
		return "String"
	case Int:
		return "Int"
	case Ident:
		return "Ident"
	case StringList:
		return "List String"
	}
	panic("unreachable")
}
