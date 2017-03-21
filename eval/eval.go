package eval

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/pkg/errors"

	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/token"
)

func New(in io.Reader, out, err io.Writer) *Evaluator {
	return &Evaluator{
		in:  in,
		out: out,
		err: err,
	}
}

func (e *Evaluator) Eval(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		err := e.eval(stmt)
		if err != nil {
			return errors.Wrap(err, "Eval")
		}
	}
	return nil
}

type Evaluator struct {
	in  io.Reader
	out io.Writer
	err io.Writer

	closeAfterStart []io.Closer
}

func (e *Evaluator) eval(stmt ast.Stmt) error {
	switch x := stmt.(type) {
	case *ast.PipeStmt:
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
		return e.execPipe(commands)
	case *ast.ExecStmt:
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

func (e *Evaluator) evalExpr(expr ast.Expr) ([]string, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		return []string{strings.Replace(x.Name, "~", os.Getenv("HOME"), -1)}, nil
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

func (e *Evaluator) execCmd(name string, args []string) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := CommandContext(ctx, name, args...)
	cmd.SetStdin(e.in)
	cmd.SetStdout(e.out)
	cmd.SetStderr(e.err)

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
}

func wait(fn func() error) <-chan error {
	c := make(chan error, 1)
	c <- fn()
	return c
}

func (e *Evaluator) execPipe(commands [][]string) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	cmds := make([]Cmd, len(commands))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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
	select {
	case s := <-c:
		return errors.New(s.String())
	case err := <-wait(f):
		return err
	}
}
