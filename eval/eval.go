package eval

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/pkg/errors"

	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/token"
)

func Eval(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		err := eval(stmt)
		if err != nil {
			return errors.Wrap(err, "Eval")
		}
	}
	return nil
}

type evaluator struct {
	in  io.Reader
	out io.Writer
	err io.Writer

	closeAfterStart []io.Closer
}

func eval(stmt ast.Stmt) error {
	switch x := stmt.(type) {
	case *ast.PipeStmt:
		e := &evaluator{
			in:  os.Stdin,
			out: os.Stdout,
			err: os.Stderr,
		}
		commands := make([][]string, 0, len(x.Args))
		for _, c := range x.Args {
			args := make([]string, 0, len(c.Args))
			for _, arg := range c.Args {
				s, err := e.evalExpr(arg)
				if err != nil {
					return err
				}
				args = append(args, s...)
			}
			if len(args) == 0 {
				return errors.New("no command to execute")
			}
			commands = append(commands, args)
		}
		return errors.Wrap(e.execPipe(commands), "eval")
	case *ast.ExecStmt:
		e := &evaluator{
			in:  os.Stdin,
			out: os.Stdout,
			err: os.Stderr,
		}
		args := make([]string, 0, len(x.Args))
		for _, arg := range x.Args {
			s, err := e.evalExpr(arg)
			if err != nil {
				return err
			}
			args = append(args, s...)
		}
		if len(args) == 0 {
			return nil
		}
		return e.execCmd(args[0], args[1:])
	}
	return fmt.Errorf("eval: unexpected type: %T", stmt)
}

func (e *evaluator) evalExpr(expr ast.Expr) ([]string, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		return []string{x.Name}, nil
	case *ast.BasicLit:
		s := strings.TrimPrefix(x.Value, "'")
		s = strings.TrimSuffix(s, "'")
		s = strings.Replace(s, "''", "'", -1)
		return []string{s}, nil
	case *ast.ParenExpr:
		var list []string
		for _, expr := range x.Exprs {
			s, err := e.evalExpr(expr)
			if err != nil {
				return nil, err
			}
			list = append(list, s...)
		}
		return list, nil
	case *ast.UnaryExpr:
		s, err := e.evalExpr(x.X)
		if err != nil {
			return nil, err
		}
		if len(s) == 0 {
			return nil, fmt.Errorf("cannot redirect")
		}
		if len(s) > 1 {
			return nil, fmt.Errorf("cannot redirect to multi-word filename")
		}
		switch x.Op {
		case token.REDIRIN:
			f, err := os.Open(s[0])
			if err != nil {
				return nil, err
			}
			e.in = f
			e.closeAfterStart = append(e.closeAfterStart, f)
		case token.REDIROUT:
			f, err := os.Create(s[0])
			if err != nil {
				return nil, err
			}
			e.out = f
			e.closeAfterStart = append(e.closeAfterStart, f)
		}
		return nil, nil
	case nil:
		return nil, nil
	}
	return nil, fmt.Errorf("evalExpr: unexpected type: %T", expr)
}

func (e *evaluator) execCmd(name string, args []string) error {
	if fn, ok := builtins[name]; ok {
		return fn(e, args)
	}
	if x, ok := aliases[name]; ok {
		name = x.cmd
		args = append(x.args, args...)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdin = e.in
	cmd.Stdout = e.out
	cmd.Stderr = e.err

	defer func() {
		for _, closer := range e.closeAfterStart {
			closer.Close()
		}
	}()
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

func (e *evaluator) execPipe(commands [][]string) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	cmds := make([]Cmd, len(commands))
	ctx, cancel := context.WithCancel(context.Background())
	for i, c := range commands {
		name := c[0]
		args := c[1:]
		cmds[i] = CommandContext(ctx, name, args...)
		if i > 0 {
			pipe, err := cmds[i-1].StdoutPipe()
			if err != nil {
				return err
			}
			cmds[i].SetStdin(pipe)
		}
		cmds[i].SetStderr(e.err)
	}
	cmds[0].SetStdin(e.in)
	cmds[len(cmds)-1].SetStdout(e.out)

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return errors.Wrap(err, "cmd.Start")
		}
	}
	f := func() error {
		for _, cmd := range cmds {
			if err := cmd.Wait(); err != nil {
				return errors.Wrap(err, "cmd.Wait")
			}
		}
		return nil
	}

	defer func() {
		for _, closer := range e.closeAfterStart {
			closer.Close()
		}
	}()
	defer cancel()
	select {
	case s := <-c:
		return errors.New(s.String())
	case err := <-wait(f):
		return err
	}
	return nil
}
