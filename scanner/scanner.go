package scanner

import (
	"fmt"
	"path/filepath"
	"unicode/utf8"

	"github.com/elpinal/coco3/token"
)

type ErrorHandler func(pos token.Position, msg string)

type Scanner struct {
	// immutable state
	file *token.File  // source file handle
	dir  string       // directory portion of file.Name()
	src  []byte       // source
	err  ErrorHandler // error reporting; or nil

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
	insertSemi bool // insert a semicolon before next newline

	// public state - ok to modify
	ErrorCount int // number of errors encountered
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

// Read the next Unicode char into s.ch.
// s.ch < 0 means end-of-file.
//
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		if s.ch == '\n' {
			s.lineOffset = s.offset
		}
		r, w := rune(s.src[s.rdOffset]), 1
		switch {
		case r == 0:
			s.error(s.offset, "illegal character NUL")
		case r >= utf8.RuneSelf:
			// not ASCII
			r, w = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError && w == 1 {
				s.error(s.offset, "illegal UTF-8 encoding")
			} else if r == bom && s.offset > 0 {
				s.error(s.offset, "illegal byte order mark")
			}
		}
		s.rdOffset += w
		s.ch = r
	} else {
		s.offset = len(s.src)
		if s.ch == '\n' {
			s.lineOffset = s.offset
		}
		s.ch = -1 // eof
	}
}

func (s *Scanner) Init(file *token.File, src []byte, err ErrorHandler) {
	// Explicitly initialize all fields since a scanner may be reused.
	if file.Size() != len(src) {
		panic(fmt.Sprintf("file size (%d) does not match src len (%d)", file.Size(), len(src)))
	}
	s.file = file
	s.dir, _ = filepath.Split(file.Name())
	s.src = src
	s.err = err

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0
	s.insertSemi = false
	s.ErrorCount = 0

	s.next()
	if s.ch == bom {
		s.next() // ignore BOM at file beginning
	}
}

func (s *Scanner) error(offs int, msg string) {
	if s.err != nil {
		s.err(s.file.Position(s.file.Pos(offs)), msg)
	}
	s.ErrorCount++
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for s.ch != ' ' && s.ch != '\t' && s.ch != '\n' && s.ch != -1 && s.ch != ';' && s.ch != '(' && s.ch != ')' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' {
		s.next()
	}
}

func (s *Scanner) Scan() (pos token.Pos, tok token.Token, lit string) {
	s.skipWhitespace()

	pos = token.Pos(s.offset)

	//insertSemi := false
	switch ch := s.ch; ch {
	default:
		lit = s.scanIdentifier()
		tok = token.IDENT
	case -1:
		s.next() // always make progress
		if s.insertSemi {
			s.insertSemi = false // EOF consumed
			return pos, token.SEMICOLON, "\n"
		}
		tok = token.EOF
	case '\n':
		s.next()             // always make progress
		s.insertSemi = false // newline consumed
		return pos, token.SEMICOLON, "\n"
	case '(':
		s.next()
		tok = token.LPAREN
	case ')':
		s.next()
		tok = token.RPAREN
	case '<':
		s.next()
		tok = token.REDIRIN
	case '>':
		s.next()
		tok = token.REDIROUT
	case ';':
		s.next()
		tok = token.SEMICOLON
		lit = ";"
	}

	return
}
