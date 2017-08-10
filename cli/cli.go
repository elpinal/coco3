package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/gate"
	"github.com/elpinal/coco3/parser"

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

	for _, alias := range c.Config.Alias {
		eval.DefAlias(alias[0], alias[1])
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
	histRunes := make([][]rune, len(history))
	for i, line := range history {
		histRunes[i] = []rune(line)
	}
	g := gate.NewContext(ctx, conf, c.In, c.Out, c.Err, histRunes)
	c.db = db
	go func(ctx context.Context) {
		for {
			if err := c.interact(g); err != nil {
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
	if err := c.execute([]byte(string(r))); err != nil {
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

var schema = `
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
