package scanner

import (
	"log"
	"unicode/utf8"

	"github.com/elpinal/coco3/token"
)

type Scanner struct {
	// immutable state
	src []byte // source

	// scanning state
	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (position after current character)
	lineOffset int  // current line offset
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

func (s *Scanner) Init(src []byte) {
	// Explicitly initialize all fields since a scanner may be reused.
	s.src = src

	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0
	s.lineOffset = 0

	s.next()
	if s.ch == bom {
		s.next() // ignore BOM at file beginning
	}
}

func (s *Scanner) error(offs int, msg string) {
	log.Println("Scanner.error:", offs, msg)
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for s.ch != ' ' && s.ch != '\t' && s.ch != -1 {
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

	switch ch := s.ch; ch {
	default:
		lit = s.scanIdentifier()
		tok = token.IDENT
	case -1:
		s.next() // always make progress
		tok = token.EOF
	}

	return
}
