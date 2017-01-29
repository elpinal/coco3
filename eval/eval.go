package eval

import (
	"errors"
	"os"
	"os/exec"

	"github.com/elpinal/coco3/ast"
)

func Eval(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		err := eval(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func eval(stmt ast.Stmt) error {
	switch x := stmt.(type) {
	case *ast.ExecStmt:
		cmdStr, err := evalExpr(x.Cmd)
		if err != nil {
			return err
		}
		args := make([]string, 0, len(x.Args))
		for _, arg := range x.Args {
			s, err := evalExpr(arg)
			if err != nil {
				return err
			}
			args = append(args, s)
		}
		return execCmd(cmdStr, args)
	}
	return errors.New("unexpected type")
}

func evalExpr(expr ast.Expr) (string, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		return x.Name, nil
	}
	return "", errors.New("unexpected type")
}

func execCmd(name string, args []string) error {
	if fn, ok := builtins[name]; ok {
		return fn(args)
	}
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
