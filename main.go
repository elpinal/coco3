package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/gate"
	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/parser"
)

type CLI struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func main() {
	c := CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(c.run(os.Args[1:]))
}

func (c CLI) run(args []string) int {
	f := flag.NewFlagSet("coco4", flag.ContinueOnError)
	f.SetOutput(c.Err)

	flagC := f.String("c", "", "take first argument as a command to execute")
	if err := f.Parse(args); err != nil {
		return 2
	}

	if *flagC != "" {
		if err := execute([]byte(*flagC)); err != nil {
			fmt.Fprintln(c.Err, err)
			return 1
		}
		return 0
	}

	if len(f.Args()) > 0 {
		for _, file := range f.Args() {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintln(c.Err, err)
				return 1
			}
			if err := execute(b); err != nil {
				fmt.Fprintln(c.Err, err)
				return 1
			}
		}
		return 0
	}

	conf := new(config.Config)
	conf.Init()
	for {
		if err := c.interact(conf); err != nil {
			fmt.Fprintln(c.Err, err)
			// return 1
		}
	}
	return 0
}

func (c CLI) interact(conf *config.Config) error {
	g := gate.New(conf, c.In, c.Out, c.Err)
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
		if err := execute([]byte(string(r))); err != nil {
			return err
		}
		g.Clear()
	}
	return nil
}

func execute(b []byte) error {
	f, err := parser.ParseSrc(b)
	if err != nil {
		return err
	}
	return eval.Eval(f.Lines)
}
