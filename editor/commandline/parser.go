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

func (s *scanner) scan() (byte, bool) {
	if s.off >= s.size {
		return 0, true
	}

	ret := s.src[s.off]

	s.off++
	return ret, false
}

func (s *scanner) lex() *token {
	ch, eof := s.scan()
	if eof {
		return nil
	}
	switch {
	case isIdent(ch):
		ch, eof := s.scan()
		if eof {
			return nil
		}
		for isIdent(ch) {
			ch, eof = s.scan()
			if eof {
				break
			}
		}
		return &token{tt: ident}
	case ch == '"':
	}
	return nil
}

func isIdent(r byte) bool {
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z'
}

type tokenType int

type token struct {
	tt tokenType
}

const (
	ident tokenType = iota
	str
)
