package eval

import (
	"context"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/pkg/errors"
)

type Cmd interface {
	Run() error
	Start() error
	Wait() error
	StderrPipe() (io.ReadCloser, error)
	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	SetStderr(io.Writer)
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetEnv([]string)
}

var _ Cmd = (*externalCmd)(nil)

func (e *Evaluator) CommandContext(ctx context.Context, name string, arg ...string) Cmd {
	if x, ok := aliases[name]; ok {
		name = x.cmd
		arg = append(x.args, arg...)
	}
	if fn, ok := builtins[name]; ok {
		return &builtinCmd{ctx: ctx, fn: fn, name: name, args: arg, e: e}
	}
	return &externalCmd{exec.Command(name, arg...)}

}

type externalCmd struct {
	*exec.Cmd
}

func (c *externalCmd) Run() error {
	err := c.Cmd.Run()
	if err == nil || isInterrupt(err) {
		return nil
	}
	return err
}

func (c *externalCmd) SetStderr(w io.Writer) {
	c.Stderr = w
}

func (c *externalCmd) SetStdin(r io.Reader) {
	c.Stdin = r
}

func (c *externalCmd) SetStdout(w io.Writer) {
	c.Stdout = w
}

func (c *externalCmd) SetEnv(env []string) {
	c.Env = env
}

type builtinCmd struct {
	ctx  context.Context
	fn   func(context.Context, stream, []string, *Evaluator, []string) error
	name string
	args []string
	e    *Evaluator
	in   io.Reader
	out  io.Writer
	err  io.Writer
	env  []string

	closeAfterStart []io.Closer
	closeAfterWait  []io.Closer
	ch              chan error
}

func (c *builtinCmd) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *builtinCmd) Start() error {
	c.ch = make(chan error)
	go func() {
		err := c.fn(c.ctx, stream{in: c.in, out: c.out, err: c.err}, c.env, c.e, c.args)
		c.ch <- err
		c.closeDescriptors(c.closeAfterStart)
	}()
	return nil
}

func (c *builtinCmd) Wait() error {
	err := <-c.ch
	c.closeDescriptors(c.closeAfterWait)
	return errors.Wrap(err, c.name)
}

func (c *builtinCmd) closeDescriptors(closers []io.Closer) {
	for _, fd := range closers {
		fd.Close()
	}
}

func (c *builtinCmd) StderrPipe() (io.ReadCloser, error) {
	if c.err != nil {
		return nil, errors.New("Stderr already set")
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	c.err = pw
	c.closeAfterStart = append(c.closeAfterStart, pw)
	c.closeAfterWait = append(c.closeAfterWait, pr)
	return pr, nil
}

func (c *builtinCmd) StdinPipe() (io.WriteCloser, error) {
	if c.in != nil {
		return nil, errors.New("Stdin already set")
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	c.in = pr
	c.closeAfterStart = append(c.closeAfterStart, pr)
	wc := &closeOnce{File: pw}
	c.closeAfterWait = append(c.closeAfterWait, wc)
	return wc, nil
}

type closeOnce struct {
	*os.File

	once sync.Once
	err  error
}

func (c *closeOnce) Close() error {
	c.once.Do(c.close)
	return c.err
}

func (c *closeOnce) close() {
	c.err = c.File.Close()
}

func (c *builtinCmd) StdoutPipe() (io.ReadCloser, error) {
	if c.out != nil {
		return nil, errors.New("Stdout already set")
	}
	pr, pw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	c.out = pw
	c.closeAfterStart = append(c.closeAfterStart, pw)
	c.closeAfterWait = append(c.closeAfterWait, pr)
	return pr, nil
}

func (c *builtinCmd) SetStderr(w io.Writer) {
	c.err = w
}

func (c *builtinCmd) SetStdin(r io.Reader) {
	c.in = r
}

func (c *builtinCmd) SetStdout(w io.Writer) {
	c.out = w
}

func (c *builtinCmd) SetEnv(env []string) {
	c.env = env
}
