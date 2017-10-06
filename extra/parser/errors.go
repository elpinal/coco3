package parser

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/elpinal/color"
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
	var buf bytes.Buffer
	buf.WriteString(p.Error())
	buf.WriteString("\n\n")

	l := fmt.Sprintf("%d: ", p.Line)
	buf.WriteString(color.Wrap(l, color.Cyan))
	buf.WriteString(strings.Split(p.Src, "\n")[p.Line-1])
	buf.WriteByte('\n')

	buf.WriteString(strings.Repeat(" ", int(p.Column-1)+len(l)))
	buf.WriteString("\033[1m^ error occurs\033[0m")
	return buf.String()
}
