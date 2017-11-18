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
	src  []byte
	size int
	off  int
}

func newScanner(src []byte) *scanner {
	return &scanner{
		src:  src,
		size: len(src),
	}
}

func (s scanner) scan() *token {
	if s.off >= s.size {
		return nil
	}

	// deal with s.src[s.off]...

	s.off++
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
