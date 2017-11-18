package commandline

import "strings"

func Parse(s string) *Command {
	ss := strings.Split(s, " ")
	if ss[0] == "" {
		return nil
	}
	return &Command{
		Name: ss[0],
		Args: ss[1:],
	}
}

type Command struct {
	Name string
	Args []string
}

type scanner struct {
	src []byte
}

func newScanner(src []byte) *scanner {
	return &scanner{
		src: src,
	}
}

func (s scanner) scan() *token {
	return nil
}

type tokenType int

type token struct {
	tt tokenType
}

const (
	ident tokenType = iota
	str
)
