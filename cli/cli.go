package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/gate"
	"github.com/elpinal/coco3/parser"

	"github.com/elpinal/coco3/extra"
	eparser "github.com/elpinal/coco3/extra/parser"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type CLI struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer

	*config.Config

	db *sqlx.DB

	execute1 func([]byte) (action, error)
}

func (c *CLI) Run(args []string) int {
	f := flag.NewFlagSet("coco3", flag.ContinueOnError)
	f.SetOutput(c.Err)
	f.Usage = func() {
		c.Err.Write([]byte("coco3 is a shell.\n"))
		c.Err.Write([]byte("Usage:\n"))
		f.PrintDefaults()
	}

	if c.Config == nil {
		c.Config = &config.Config{}
	}
	c.Config.Init()
	flagC := f.String("c", "", "take first argument as a command to execute")
	flagE := f.Bool("extra", c.Config.Extra, "switch to extra mode")
	if err := f.Parse(args); err != nil {
		return 2
	}
	return c.run(f.Args(), flagC, flagE)
}

func (c *CLI) run(args []string, flagC *string, flagE *bool) int {
	for _, alias := range c.Config.Alias {
		eval.DefAlias(alias[0], alias[1])
	}

	for k, v := range c.Config.Env {
		err := os.Setenv(k, v)
		if err != nil {
			c.errorln(err)
			return 1
		}
	}

	setpath(c.Config.Paths)

	if *flagE {
		// If -extra flag is on, enable extra mode on any command executions.
		c.execute1 = c.executeExtra
	} else {
		c.execute1 = c.execute
	}

	if len(c.Config.StartUpCommand) > 0 {
		a, err := c.execute1(c.Config.StartUpCommand)
		if err != nil {
			c.printExecError(err)
			return 1
		}
		if e, ok := a.(exit); ok {
			return e.code
		}
	}

	if *flagC != "" {
		a, err := c.execute1([]byte(*flagC))
		if err != nil {
			c.printExecError(err)
			return 1
		}
		if e, ok := a.(exit); ok {
			return e.code
		}
		return 0
	}

	if len(args) > 0 {
		a, err := c.runFiles(args)
		if err != nil {
			c.printExecError(err)
			return 1
		}
		if e, ok := a.(exit); ok {
			return e.code
		}
		return 0
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	histRunes, err := c.getHistory(c.Config.HistFile)
	if err != nil {
		c.errorln(err)
		return 1
	}
	g := gate.NewContext(ctx, c.Config, c.In, c.Out, c.Err, histRunes)
	for {
		a, err := c.interact(g)
		if err != nil {
			c.printExecError(err)
		}
		if e, ok := a.(exit); ok {
			return e.code
		}
	}
	return 0
}

func (c *CLI) errorf(s string, a ...interface{}) {
	fmt.Fprintf(c.Err, s, a...)
}

func (c *CLI) errorln(s ...interface{}) {
	fmt.Fprintln(c.Err, s...)
}

func (c *CLI) errorp(s ...interface{}) {
	fmt.Fprint(c.Err, s...)
}

func (c *CLI) getHistory(filename string) ([][]rune, error) {
	db, err := sqlx.Connect("sqlite3", filename)
	if err != nil {
		return nil, errors.Wrap(err, "connecting history file")
	}
	_, err = db.Exec(schema)
	if err != nil {
		return nil, errors.Wrap(err, "initializing history file")
	}
	var history []string
	err = db.Select(&history, "select line from command_info")
	if err != nil {
		return nil, errors.Wrap(err, "restoring history")
	}
	// TODO: Is this way proper?
	c.db = db
	return sanitizeHistory(history), nil
}

func (c *CLI) printExecError(err error) {
	if pe, ok := err.(*eparser.ParseError); ok {
		c.errorln(pe.Verbose())
	} else {
		c.errorln(err)
	}
}

// setpath sets the PATH environment variable.
func setpath(args []string) {
	if len(args) == 0 {
		return
	}
	paths := filepath.SplitList(os.Getenv("PATH"))
	var newPaths []string
	for _, path := range paths {
		if contains(args, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	newPaths = append(args, newPaths...)
	os.Setenv("PATH", strings.Join(newPaths, string(filepath.ListSeparator)))
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func sanitizeHistory(history []string) [][]rune {
	histRunes := make([][]rune, 0, len(history))
	for _, line := range history {
		if line == "" {
			continue
		}
		l := len(histRunes)
		s := []rune(line)
		if l > 0 && compareRunes(histRunes[l-1], s) {
			continue
		}
		histRunes = append(histRunes, s)
	}
	return histRunes
}

func compareRunes(r1, r2 []rune) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, r := range r1 {
		if r2[i] != r {
			return false
		}
	}
	return true
}

func (c *CLI) interact(g gate.Gate) (action, error) {
	r, end, err := c.read(g)
	if err != nil {
		return nil, err
	}
	if end {
		return exit{code: 0}, nil
	}
	ch := c.writeHistory(r)
	a, err := c.execute1([]byte(string(r)))
	if err != nil {
		return a, err
	}
	return a, <-ch
}

func (c *CLI) read(g gate.Gate) ([]rune, bool, error) {
	defer c.Out.Write([]byte{'\n'})
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		return nil, false, err
	}
	defer func() {
		if err := terminal.Restore(0, oldState); err != nil {
			c.errorln(err)
		}
	}()
	r, end, err := g.Read()
	if err != nil {
		return nil, false, err
	}
	return r, end, nil
}

func (c *CLI) writeHistory(r []rune) <-chan error {
	startTime := time.Now()
	ch := make(chan error)
	go func() {
		_, err := c.db.Exec("insert into command_info (time, line) values ($1, $2)", startTime, string(r))
		if err != nil {
			ch <- errors.Wrap(err, "saving history")
			return
		}
		ch <- nil
	}()
	return ch
}

const schema = `
create table if not exists command_info (
    time datetime,
    line text
)`

func (c *CLI) execute(b []byte) (action, error) {
	f, err := parser.ParseSrc(b)
	if err != nil {
		return nil, err
	}
	e := eval.New(c.In, c.Out, c.Err, c.db)
	err = e.Eval(f.Lines)
	select {
	case code := <-e.ExitCh:
		return exit{code: code}, nil
	default:
	}
	return nil, err
}

func (c *CLI) executeExtra(b []byte) (action, error) {
	cmd, err := eparser.Parse(b)
	if err != nil {
		return nil, err
	}
	e := extra.New(extra.Option{DB: c.db})
	err = e.Eval(cmd)
	if err == nil {
		return nil, nil
	}
	if pe, ok := err.(*eparser.ParseError); ok {
		pe.Src = string(b)
	}
	return nil, err
}

func (c *CLI) runFiles(files []string) (action, error) {
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		a, err := c.execute1(b)
		if err != nil {
			return nil, err
		}
		if a != nil {
			return a, nil
		}
	}
	return nil, nil
}

type action interface {
	act()
}

type exit struct {
	code int
}

func (e exit) act() {}
