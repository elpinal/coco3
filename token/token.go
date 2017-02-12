package token

import "strconv"

// Token is the set of lexical tokens of the coco3.
type Token int

// The list of tokens.
const (
	// Special tokens
	ILLEGAL Token = iota
	EOF

	literal_beg
	IDENT // main
	literal_end

	operator_beg
	LPAREN // (
	RPAREN // )

	SEMICOLON // ;
	operator_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",

	EOF: "EOF",

	IDENT: "IDENT",

	LPAREN: "(",
	RPAREN: ")",

	SEMICOLON: ";",
}

func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func (tok Token) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

func (tok Token) IsOperator() bool { return operator_beg < tok && tok < operator_end }
