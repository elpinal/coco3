package main

import (
	"os"

	"github.com/elpinal/coco3/cli"
)

func main() {
	c := cli.CLI{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
	os.Exit(c.Run(os.Args[1:]))
}
