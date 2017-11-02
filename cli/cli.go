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

	config.Config

	db *sqlx.DB

	exitCh chan int
	doneCh chan struct{} // to ensure exiting just after exitCh received

	execute1 func([]byte) error
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
	flagE := f.Bool("extra", c.Config.Extra, "switch to extra mode")
	if err := f.Parse(args); err != nil {
		return 2
	}

	for _, alias := range c.Config.Alias {
		eval.DefAlias(alias[0], alias[1])
	}

	for k, v := range c.Config.Env {
		err := os.Setenv(k, v)
		if err != nil {
			fmt.Fprintln(c.Err, err)
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
		done := make(chan struct{})
		go func() {
			if err := c.execute1(c.Config.StartUpCommand); err != nil {
				c.printExecError(err)
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
			if err := c.execute1([]byte(*flagC)); err != nil {
				c.printExecError(err)
				c.exitCh <- 1
				return
			}
			c.exitCh <- 0
		}()
		return <-c.exitCh
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if len(f.Args()) > 0 {
		go c.runFiles(ctx, f.Args())
		return <-c.exitCh
	}

	conf := &c.Config
	conf.Init()
	db, err := sqlx.Connect("sqlite3", conf.HistFile)
	if err != nil {
		fmt.Fprintf(c.Err, "connecting history file: %v\n", err)
		return 1
	}
	_, err = db.Exec(schema)
	if err != nil {
		fmt.Fprintf(c.Err, "initializing history file: %v\n", err)
		return 1
	}
	var history []string
	err = db.Select(&history, "select line from command_info")
	if err != nil {
		fmt.Fprintf(c.Err, "restoring history: %v\n", err)
		return 1
	}
	histRunes := sanitizeHistory(history)
	g := gate.NewContext(ctx, conf, c.In, c.Out, c.Err, histRunes)
	c.db = db
	go func(ctx context.Context) {
		for {
			if err := c.interact(g); err != nil {
				c.printExecError(err)
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

func (c *CLI) printExecError(err error) {
	if pe, ok := err.(*eparser.ParseError); ok {
		fmt.Fprintln(c.Err, pe.Verbose())
	} else {
		fmt.Fprintln(c.Err, err)
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

func (c *CLI) interact(g gate.Gate) error {
	r, end, err := c.read(g)
	if err != nil {
		return err
	}
	if end {
		c.exitCh <- 0
		<-c.doneCh
		return nil
	}
	go c.writeHistory(r)
	if err := c.execute1([]byte(string(r))); err != nil {
		return err
	}
	g.Clear()
	return nil
}

func (c *CLI) read(g gate.Gate) ([]rune, bool, error) {
	defer c.Out.Write([]byte{'\n'})
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := terminal.Restore(0, oldState); err != nil {
			fmt.Fprintln(c.Err, err)
		}
	}()
	r, end, err := g.Read()
	if err != nil {
		return nil, false, err
	}
	return r, end, nil
}

func (c *CLI) writeHistory(r []rune) {
	startTime := time.Now()
	_, err := c.db.Exec("insert into command_info (time, line) values ($1, $2)", startTime, string(r))
	if err != nil {
		fmt.Fprintf(c.Err, "saving history: %v\n", err)
		c.exitCh <- 1
	}
}

const schema = `
create table if not exists command_info (
    time datetime,
    line text
)`

func (c *CLI) execute(b []byte) error {
	f, err := parser.ParseSrc(b)
	if err != nil {
		return err
	}
	e := eval.New(c.In, c.Out, c.Err, c.db)
	err = e.Eval(f.Lines)
	select {
	case code := <-e.ExitCh:
		c.exitCh <- code
		<-c.doneCh
	default:
	}
	return err
}

func (c *CLI) executeExtra(b []byte) error {
	cmd, err := eparser.Parse(b)
	if err != nil {
		return err
	}
	e := extra.New(extra.Option{DB: c.db})
	err = e.Eval(cmd)
	if err == nil {
		return nil
	}
	if pe, ok := err.(*eparser.ParseError); ok {
		pe.Src = string(b)
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
		if err := c.execute1(b); err != nil {
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
