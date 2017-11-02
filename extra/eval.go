package extra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"

	"github.com/elpinal/coco3/extra/ast"
	"github.com/elpinal/coco3/extra/parser" // Only for ParseError.
	"github.com/elpinal/coco3/extra/typed"
	"github.com/elpinal/coco3/extra/types"
)

type Env struct {
	cmds map[string]typed.Command
	Option
}

type Option struct {
	DB *sqlx.DB
}

func New(opt Option) Env {
	return Env{
		Option: opt,
		cmds: map[string]typed.Command{
			"exec":    execCommand,
			"cd":      cdCommand,
			"exit":    exitCommand,
			"free":    freeCommand,
			"history": historyCommand,

			"ls":  lsCommand,
			"man": manCommand,

			"git":   gitCommand,
			"cargo": cargoCommand,
			"go":    goCommand,
			"stack": stackCommand,

			"vim":    vimCommand,
			"emacs":  emacsCommand,
			"screen": screenCommand,
		},
	}
}

func WithoutDefault() Env {
	return Env{cmds: make(map[string]typed.Command)}
}

func (e *Env) Bind(name string, c typed.Command) {
	e.cmds[name] = c
}

func (e *Env) Eval(command *ast.Command) error {
	if command == nil {
		return nil
	}
	tc, found := e.cmds[command.Name.Lit]
	if !found {
		return &parser.ParseError{
			Msg:    fmt.Sprintf("no such typed command: %q", command.Name.Lit),
			Line:   command.Name.Line,
			Column: command.Name.Column,
		}
	}
	if len(command.Args) != len(tc.Params) {
		return &parser.ParseError{
			Msg:    fmt.Sprintf("the length of args (%d) != the one of params (%d)", len(command.Args), len(tc.Params)),
			Line:   command.Name.Line,
			Column: command.Name.Column,
		}
	}
	for i, arg := range command.Args {
		if arg.Type() != tc.Params[i] {
			return &parser.ParseError{
				Msg:    fmt.Sprintf("type mismatch: (%v) (type of %v) does not match with (%v) (expected type)", arg.Type(), arg, tc.Params[i]),
				Line:   command.Name.Line,
				Column: command.Name.Column,
			}
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer close(c)
	defer signal.Stop(c)

	return tc.Fn(command.Args, e.DB)
}

func toSlice(list ast.List) ([]string, error) {
	ret := make([]string, 0, list.Length())
	for {
		switch x := list.(type) {
		case *ast.Cons:
			ret = append(ret, x.Head)
			list = x.Tail
		case *ast.Empty:
			return ret, nil
		default:
			return nil, fmt.Errorf("unexpected list type: %T", x)
		}
	}
}

var execCommand = typed.Command{
	Params: []types.Type{types.String, types.StringList},
	Fn: func(args []ast.Expr, _ *sqlx.DB) error {
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, "exec")
		}
		cmd := stdCmd(args[0].(*ast.String).Lit, cmdArgs...)
		return cmd.Run()
	},
}

var cdCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn: func(args []ast.Expr, _ *sqlx.DB) error {
		return os.Chdir(args[0].(*ast.String).Lit)
	},
}

var exitCommand = typed.Command{
	Params: []types.Type{types.Int},
	Fn: func(args []ast.Expr, _ *sqlx.DB) error {
		n, err := strconv.Atoi(args[0].(*ast.Int).Lit)
		if err != nil {
			return err
		}
		os.Exit(n)
		return nil
	},
}

var freeCommand = typed.Command{
	Params: []types.Type{types.String, types.StringList},
	Fn: func(args []ast.Expr, _ *sqlx.DB) error {
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, "free")
		}
		name := args[0].(*ast.String).Lit
		cmd := exec.Cmd{Path: name, Args: append([]string{name}, cmdArgs...)}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	},
}

func commandsInCommand(name string) func([]ast.Expr, *sqlx.DB) error {
	return func(args []ast.Expr, _ *sqlx.DB) error {
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, name)
		}
		var cmd *exec.Cmd
		switch lit := args[0].(*ast.Ident).Lit; lit {
		case "command":
			cmd = stdCmd(name, cmdArgs...)
		default:
			cmd = stdCmd(name, append([]string{lit}, cmdArgs...)...)
		}
		return cmd.Run()
	}
}

var gitCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("git"),
}

var cargoCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("cargo"),
}

var goCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("go"),
}

var stackCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("stack"),
}

func stdCmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func stdExec(name string, args ...string) func([]ast.Expr, *sqlx.DB) error {
	return func(_ []ast.Expr, _ *sqlx.DB) error {
		return stdCmd(name, args...).Run()
	}
}

var vimCommand = typed.Command{
	Params: []types.Type{},
	Fn:     stdExec("vim"),
}

var emacsCommand = typed.Command{
	Params: []types.Type{},
	Fn:     stdExec("emacs"),
}

func withEnv(s string, cmd *exec.Cmd) *exec.Cmd {
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}
	cmd.Env = append(cmd.Env, s)
	return cmd
}

var screenCommand = typed.Command{
	Params: []types.Type{},
	Fn: func(_ []ast.Expr, _ *sqlx.DB) error {
		return withEnv("LANG=en_US.UTF-8", stdCmd("screen")).Run()
	},
}

type execution struct {
	Time time.Time
	Line string
}

var historyCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn: func(e []ast.Expr, db *sqlx.DB) error {
		var jsonFormat bool
		var enc *json.Encoder
		switch format := e[0].(*ast.String).Lit; format {
		case "json":
			jsonFormat = true
		case "lines":
		default:
			return fmt.Errorf("history: format %q is not supported", format)
		}

		buf := bufio.NewWriter(os.Stdout)
		if jsonFormat {
			enc = json.NewEncoder(buf)
		}
		rows, err := db.Queryx("select * from command_info")
		if err != nil {
			return err
		}
		data := execution{}
		for rows.Next() {
			err := rows.StructScan(&data)
			if err != nil {
				return err
			}
			if jsonFormat {
				err := enc.Encode(data)
				if err != nil {
					return err
				}
			} else {
				buf.WriteString(data.Time.Format("Mon, 02 Jan 2006 15:04:05"))
				buf.Write([]byte("  "))
				buf.WriteString(data.Line)
				buf.WriteByte('\n')
			}
		}
		return buf.Flush()
	},
}

var lsCommand = typed.Command{
	Params: []types.Type{},
	Fn:     stdExec("ls", "--show-control-chars", "--color=auto"),
}

var manCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn: func(e []ast.Expr, _ *sqlx.DB) error {
		lit := e[0].(*ast.String).Lit
		return stdCmd("man", lit).Run()
	},
}
