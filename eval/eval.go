package eval

import (
	"context"
	"errors"
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
