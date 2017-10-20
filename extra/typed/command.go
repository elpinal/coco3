package typed

import (
	"github.com/jmoiron/sqlx"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/types"
)

type Command struct {
	Params []types.Type
	Fn     func([]ast.Expr, *sqlx.DB) error
}
