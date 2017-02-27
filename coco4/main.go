package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/parser"
)

type cli struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

func main() {
	c := cli{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
	os.Exit(c.run(os.Args[1:]))
}

func (c cli) run(args []string) int {
	f := flag.NewFlagSet("coco4", flag.ContinueOnError)
	f.SetOutput(c.err)

	flagC := f.String("c", "", "take first argument as a command to execute")
	if err := f.Parse(args); err != nil {
		return 2
	}

	if *flagC != "" {
		if err := execute([]byte(*flagC)); err != nil {
			fmt.Fprintln(c.err, err)
			return 1
		}
		return 0
	}

	if len(f.Args()) > 0 {
		for _, file := range f.Args() {
			b, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Fprintln(c.err, err)
				return 1
			}
			if err := execute(b); err != nil {
				fmt.Fprintln(c.err, err)
				return 1
			}
		}
		return 0
	}

	if err := interact(); err != nil {
		fmt.Fprintln(c.err, err)
		return 1
	}
	return 0
}

func interact() error {
	for {
		// TODO
		return nil
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
