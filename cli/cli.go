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
	doneCh chan struct{}
}

func (c *CLI) Run(args []string) int {
	c.exitCh = make(chan int)
	c.doneCh = make(chan struct{})
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
		go func() {
			if err := c.execute(c.Config.StartUpCommand); err != nil {
				fmt.Fprintln(c.Err, err)
				c.exitCh <- 1
			}
			c.exitCh <- 0
		}()
		i := <-c.exitCh
		return i
	}

	if *flagC != "" {
		go func() {
			if err := c.execute([]byte(*flagC)); err != nil {
				fmt.Fprintln(c.Err, err)
				c.exitCh <- 1
			}
			c.exitCh <- 0
		}()
		i := <-c.exitCh
		return i
	}

	if len(f.Args()) > 0 {
		defer func() {
			close(c.doneCh)
		}()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func(ctx context.Context) {
			for _, file := range f.Args() {
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
		}(ctx)
		i := <-c.exitCh
		return i
	}

	conf := &c.Config
	conf.Init()
	g := gate.New(conf, c.In, c.Out, c.Err)
	go func() {
		for {
			if err := c.interact(g); err != nil {
				fmt.Fprintln(c.Err, err)
				g.Clear()
			}
		}
	}()
	i := <-c.exitCh
	return i
}

func (c *CLI) interact(g gate.Gate) error {
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
