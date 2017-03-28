package cli

import (
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
}

func (c *CLI) Run(args []string) (code int) {
	c.exitCh = make(chan int, 1)
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

	defer func() {
		select {
		case i := <-c.exitCh:
			code = i
		default:
		}
	}()

	if len(c.Config.StartUpCommand) > 0 {
		if err := c.execute(c.Config.StartUpCommand); err != nil {
			fmt.Fprintln(c.Err, err)
			return 1
		}
		select {
		case i := <-c.exitCh:
			return i
		default:
		}
	}

	if *flagC != "" {
		if err := c.execute([]byte(*flagC)); err != nil {
			fmt.Fprintln(c.Err, err)
			return 1
		}
		return code
	}

	if len(f.Args()) > 0 {
		for _, file := range f.Args() {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintln(c.Err, err)
				return 1
			}
			if err := c.execute(b); err != nil {
				fmt.Fprintln(c.Err, err)
				return 1
			}
		}
		return code
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
	default:
	}
	return err
}
