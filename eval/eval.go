package eval

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

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
		list, err := evalExpr(x.Cmd)
		if err != nil {
			return err
		}
		cmdStr := list[0]
		list = list[1:]
		args := make([]string, len(list), len(x.Args)+len(list))
		copy(args, list)
		for _, arg := range x.Args {
			s, err := evalExpr(arg)
			if err != nil {
				return err
			}
			args = append(args, s...)
		}
		return execCmd(cmdStr, args)
	}
	return fmt.Errorf("eval: unexpected type: %T", stmt)
}

func evalExpr(expr ast.Expr) ([]string, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		return []string{x.Name}, nil
	case *ast.ParenExpr:
		var list []string
		for _, e := range x.Exprs {
			s, err := evalExpr(e)
			if err != nil {
				return nil, err
			}
			list = append(list, s...)
		}
		return list, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("evalExpr: unexpected type: %T", expr)
}

func execCmd(name string, args []string) error {
	if fn, ok := builtins[name]; ok {
		return fn(args)
	}
	if x, ok := aliases[name]; ok {
		name = x.cmd
		args = append(x.args, args...)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	defer cancel()
	select {
	case s := <-c:
		return errors.New(s.String())
	case err := <-wait(cmd.Run):
		return err
	}
	return nil
}

func wait(fn func() error) <-chan error {
	c := make(chan error, 1)
	c <- fn()
	return c
}
