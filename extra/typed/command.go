package typed

import (
	"bytes"

	"github.com/jmoiron/sqlx"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/types"
)

type Command struct {
	Params []types.Type
	Fn     func([]ast.Expr, *sqlx.DB) error
}

func (c *Command) Signature() []byte {
	if len(c.Params) == 0 {
		return nil
	}
	var buf bytes.Buffer
	buf.WriteString(c.Params[0].String())
	if len(c.Params) == 1 {
		return buf.Bytes()
	}
	for _, t := range c.Params[1:] {
		buf.Write([]byte(" -> "))
		buf.WriteString(t.String())
	}
	return buf.Bytes()
}
