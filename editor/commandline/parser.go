package commandline

import (
	"fmt"
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
	src   []byte
	size  int
	start int
	off   int

	tokens chan token
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

func (s *scanner) scan() (*token, error) {
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
				return nil, s.error("unexpected eof in string literal")
			}
			ret = append(ret, ch)
		}
		return &token{tt: str, value: ret}, nil
	case isWhitespace(ch):
		return s.scan()
	}
	return nil, s.error("unexpected character")
}

func (s *scanner) error(msg string) *scanError {
	return &scanError{
		msg: msg,
		off: s.off,
	}
}

type scanError struct {
	msg string
	off int
}

func (s *scanError) Error() string {
	return fmt.Sprintf("error at offset %d: %s", s.off, s.msg)
}

func isWhitespace(b byte) bool {
	return b == ' '
}

func isIdent(b byte) bool {
	return 'A' <= b && b <= 'Z' || 'a' <= b && b <= 'z'
}

type stateFn func(*scanner) stateFn

func scan(src []byte) *scanner {
	s := &scanner{
		src:    src,
		size:   len(src),
		tokens: make(chan token),
	}
	go s.run()
	return s
}

func (s *scanner) run() {
	for state := scanToken; state != nil; {
		state = state(s)
	}
	close(s.tokens)
}

func scanToken(s *scanner) stateFn {
	for {
		b, eof := s.next()
		if eof {
			break
		}
		switch {
		case isWhitespace(b):
			s.start = s.off
		case isIdent(b):
			return scanIdent
		case b == '"':
			return scanString
		default:
			s.emit(tokenErr)
			return nil
		}
	}
	s.emit(tokenEOF)
	return nil
}

func scanIdent(s *scanner) stateFn {
	for {
		b, eof := s.next()
		if eof {
			break
		}
		if !isIdent(b) {
			s.off--
			break
		}
	}
	if s.start < s.off {
		s.emit(ident)
		return scanToken
	}
	s.emit(tokenEOF)
	return nil
}

func scanString(s *scanner) stateFn {
	for {
		b, eof := s.next()
		if eof {
			break
		}
		if b == '"' {
			s.emit(str) // including quotes
			return scanToken
		}
	}
	if s.start < s.off {
		s.emit(tokenErr)
		return nil
	}
	s.emit(tokenEOF)
	return nil
}

func (s *scanner) emit(tt tokenType) {
	s.tokens <- token{
		tt:    tt,
		value: s.src[s.start:s.off],
	}
	s.start = s.off
}

type tokenType int

type token struct {
	tt    tokenType
	value []byte
}

const (
	tokenErr tokenType = iota
	tokenEOF
	ident
	str
)
