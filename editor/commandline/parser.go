package commandline

import (
	"errors"
	"strings"
)

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

func (s *scanner) next() (byte, bool) {
	if s.off >= s.size {
		return 0, true
	}

	ret := s.src[s.off]

	s.off++
	return ret, false
}

func (s *scanner) lex() (*token, error) {
	ch, eof := s.next()
	if eof {
		return nil, nil
	}
	switch {
	case isIdent(ch):
		var ret []byte
		for isIdent(ch) {
			ch, eof = s.next()
			if eof {
				break
			}
			ret = append(ret, ch)
		}
		return &token{tt: ident, value: ret}, nil
	case ch == '"':
		var ret []byte
		for ch != '"' {
			ch, eof = s.next()
			if eof {
				return nil, errors.New("unexpected eof in string literal")
			}
			ret = append(ret, ch)
		}
		return &token{tt: str, value: ret}, nil
	}
	return nil, nil
}

func isIdent(r byte) bool {
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z'
}

type tokenType int

type token struct {
	tt    tokenType
	value []byte
}

const (
	ident tokenType = iota
	str
)
