package types

type Type int

const (
	String Type = iota + 1
	Ident
	StringList
)
