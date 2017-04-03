package eval

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	"github.com/elpinal/coco3/ast"
	"github.com/elpinal/coco3/token"
)

var ErrInterrupted = errors.New("signal caught: interrupt")

func New(in io.Reader, out, err io.Writer) *Evaluator {
	return &Evaluator{
		in:     in,
		out:    out,
		err:    err,
		ExitCh: make(chan int, 1),
	}
}

func (e *Evaluator) Eval(stmts []ast.Stmt) error {
	for _, stmt := range stmts {
		err := e.eval(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

type Evaluator struct {
	in  io.Reader
	out io.Writer
	err io.Writer

	closeAfterStart []io.Closer

	ExitCh chan int
}

func (e *Evaluator) eval(stmt ast.Stmt) error {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer close(c)
	go func() {
		<-c
		cancel()
	}()
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
		return e.execPipe(ctx, commands)
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
		return e.execCmd(ctx, args[0], args[1:])
	}
	return fmt.Errorf("unexpected type: %T", stmt)
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
	return nil, fmt.Errorf("unexpected type: %T", expr)
}

func (e *Evaluator) execCmd(ctx context.Context, name string, args []string) error {
	cmd := e.CommandContext(ctx, name, args...)
	cmd.SetStdin(e.in)
	cmd.SetStdout(e.out)
	cmd.SetStderr(e.err)
	return e.run(cmd)
}

func wait(fn func() error) <-chan error {
	c := make(chan error)
	go func() {
		c <- fn()
	}()
	return c
}

func (e *Evaluator) execPipe(ctx context.Context, commands [][]string) error {
	cmds, err := e.makePipe(ctx, commands)
	if err != nil {
		return errors.Wrap(err, "constructing pipe")
	}
	return e.run(cmds)
}

type runner interface {
	Run() error
}

func (e *Evaluator) run(cmd runner) error {
	defer func() {
		for _, closer := range e.closeAfterStart {
			closer.Close()
		}
	}()
	return <-wait(cmd.Run)
}

func (e *Evaluator) makePipe(ctx context.Context, commands [][]string) (pipeCmd, error) {
	cmds := make([]Cmd, len(commands))
	for i, c := range commands {
		name := c[0]
		args := c[1:]
		cmds[i] = e.CommandContext(ctx, name, args...)
		if i > 0 {
			pipe, err := cmds[i-1].StdoutPipe()
			if err != nil {
				return pipeCmd{}, err
			}
			cmds[i].SetStdin(pipe)
		}
		cmds[i].SetStderr(e.err)
	}
	cmds[0].SetStdin(e.in)
	cmds[len(cmds)-1].SetStdout(e.out)
	return pipeCmd{cmds: cmds, errStream: e.err}, nil
}

type pipeCmd struct {
	cmds      []Cmd
	errStream io.Writer
}

func (p pipeCmd) errorf(format string, err error) {
	fmt.Fprintf(p.errStream, format, err)
}

func (p pipeCmd) start() error {
	for _, cmd := range p.cmds {
		if err := cmd.Start(); err != nil {
			p.errorf("%v\n", err)
		}
	}
	return nil
}

func isInterrupt(err error) bool {
	ee, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}
	status := ee.ProcessState.Sys().(syscall.WaitStatus)
	return status.Signal() == syscall.SIGINT
}

func (p pipeCmd) wait() error {
	for _, cmd := range p.cmds {
		err := cmd.Wait()
		if err == nil || isInterrupt(err) {
			continue
		}
		p.errorf("%v\n", err)
	}
	return nil
}

func (p pipeCmd) Run() error {
	err := p.start()
	if err != nil {
		return err
	}
	return p.wait()
}
