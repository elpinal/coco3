package commandline

import (
	"errors"
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

func ParseT(s string) (*CommandT, error) {
	return parse([]byte(s))
}

type CommandT struct {
	Name []byte
	Args []Token
}

func parse(src []byte) (*CommandT, error) {
	s := scan(src)
	id, err := parseIdent(s.tokens)
	if err != nil {
		return nil, err
	}
	var args []Token
	for t := range s.tokens {
		if t.Type == tokenEOF {
			break
		}
		args = append(args, t)
	}
	return &CommandT{Name: id.Value, Args: args}, nil
}

func parseIdent(ch chan Token) (Token, error) {
	t := <-ch
	if t.Type == tokenEOF {
		return t, nil
	}
	if t.Type != tokenIdent {
		return Token{}, errors.New("not identifier")
	}
	return t, nil
}

type scanner struct {
	src   []byte
	size  int
	start int
	off   int

	tokens chan Token
}

func (s *scanner) next() (byte, bool) {
	if s.off >= s.size {
		return 0, true
	}

	ret := s.src[s.off]

	s.off++
	return ret, false
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
	return fmt.Sprintf("byte offset %d: %s", s.off, s.msg)
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
		tokens: make(chan Token),
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
		s.emit(tokenIdent)
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
			s.emit(tokenString) // including quotes
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

func (s *scanner) emit(t tokenType) {
	s.tokens <- Token{
		Type:  t,
		Value: s.src[s.start:s.off],
	}
	s.start = s.off
}

type tokenType int

type Token struct {
	Type  tokenType
	Value []byte
}

const (
	tokenErr tokenType = iota
	tokenEOF
	tokenIdent
	tokenString
)
