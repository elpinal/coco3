package parser

import (
	"fmt"
	"strings"
)

type ParseError struct {
	// starts with 1
	Line   uint
	Column uint

	Msg string

	Src string
}

func (p *ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s", p.Line, p.Column, p.Msg)
}

func (p *ParseError) Verbose() string {
	l := fmt.Sprintf("%d: ", p.Line)
	return p.Error() + "\n\n" +
		"\033[36m" + l + "\033[0m" + strings.Split(p.Src, "\n")[p.Line-1] + "\n" +
		strings.Repeat(" ", int(p.Column-1)+len(l)) + "\033[1m" + "^ error occurs" + "\033[0m"
}
