package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/gate"
	"github.com/elpinal/coco3/parser"
)

type CLI struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	config.Config

	exitCh chan int
	doneCh chan struct{} // to ensure exiting just after exitCh received
}

func (c *CLI) Run(args []string) int {
	c.exitCh = make(chan int)
	c.doneCh = make(chan struct{})
	defer func() {
		close(c.doneCh)
	}()

	f := flag.NewFlagSet("coco3", flag.ContinueOnError)
	f.SetOutput(c.Err)
	f.Usage = func() {
		c.Err.Write([]byte("coco3 is a shell.\n"))
		c.Err.Write([]byte("Usage:\n"))
		f.PrintDefaults()
	}

	flagC := f.String("c", "", "take first argument as a command to execute")
	if err := f.Parse(args); err != nil {
		return 2
	}

	if len(c.Config.StartUpCommand) > 0 {
		done := make(chan struct{})
		go func() {
			if err := c.execute(c.Config.StartUpCommand); err != nil {
				fmt.Fprintln(c.Err, err)
				c.exitCh <- 1
			}
			close(done)
		}()
		select {
		case code := <-c.exitCh:
			return code
		case <-done:
		}
	}

	if *flagC != "" {
		go func() {
			if err := c.execute([]byte(*flagC)); err != nil {
				fmt.Fprintln(c.Err, err)
				c.exitCh <- 1
				return
			}
			c.exitCh <- 0
		}()
		return <-c.exitCh
	}

	if len(f.Args()) > 0 {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go c.runFiles(ctx, f.Args())
		return <-c.exitCh
	}

	conf := &c.Config
	conf.Init()
	g := gate.New(conf, c.In, c.Out, c.Err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
		for {
			if err := c.interact(ctx, g); err != nil {
				fmt.Fprintln(c.Err, err)
				g.Clear()
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ctx)
	return <-c.exitCh
}

func (c *CLI) interact(ctx context.Context, g gate.Gate) error {
	for {
		old, err := enterRowMode()
		if err != nil {
			return err
		}
		r, err := g.Read()
		if err != nil {
			return err
		}
		c.Out.Write([]byte{'\n'})
		if err := exitRowMode(old); err != nil {
			return err
		}
		if err := c.execute([]byte(string(r))); err != nil {
			return err
		}
		g.Clear()
		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (c *CLI) execute(b []byte) error {
	f, err := parser.ParseSrc(b)
	if err != nil {
		return err
	}
	e := eval.New(c.In, c.Out, c.Err)
	err = e.Eval(f.Lines)
	select {
	case code := <-e.ExitCh:
		c.exitCh <- code
		<-c.doneCh
	default:
	}
	return err
}

func (c *CLI) runFiles(ctx context.Context, files []string) {
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Fprintln(c.Err, err)
			c.exitCh <- 1
			return
		}
		if err := c.execute(b); err != nil {
			fmt.Fprintln(c.Err, err)
			c.exitCh <- 1
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
		}
	}
	c.exitCh <- 0
}
