package extra

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
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
			"execenv": execenvCommand, // exec with env
			"cd":      cdCommand,
			"exit":    exitCommand,
			"free":    freeCommand,
			"history": historyCommand,

			"remove": removeCommand,

			"ls":   lsCommand,
			"man":  manCommand,
			"make": makeCommand,

			"git":    gitCommand,
			"rustup": rustupCommand,
			"cargo":  cargoCommand,
			"go":     goCommand,
			"stack":  stackCommand,
			"lein":   leinCommand,
			"ocaml":  ocamlCommand,

			"vim":    vimCommand,
			"emacs":  emacsCommand,
			"screen": screenCommand,

			"cnp":  cnpCommand,
			"gvmn": gvmnCommand,
			"vvmn": vvmnCommand,
		},
	}
}

func WithoutDefault() Env {
	return Env{cmds: make(map[string]typed.Command)}
}

func (e *Env) Bind(name string, c typed.Command) {
	e.cmds[name] = c
}

func (e *Env) Eval(command *ast.Command) (err error) {
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

	defer func() {
		r := recover()
		if r == nil {
			return
		}
		var ok bool
		// may overwrite error of tc.Fn.
		err, ok = r.(error)
		if !ok {
			panic(r)
		}
	}()

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

var execenvCommand = typed.Command{
	Params: []types.Type{types.StringList, types.String, types.StringList},
	Fn: func(args []ast.Expr, _ *sqlx.DB) error {
		cmdArgs, err := toSlice(args[2].(ast.List))
		if err != nil {
			return errors.Wrap(err, "execenv")
		}
		cmd := stdCmd(args[1].(*ast.String).Lit, cmdArgs...)
		env, err := toSlice(args[0].(ast.List))
		if err != nil {
			return errors.Wrap(err, "execenv")
		}
		for _, e := range env {
			if !strings.Contains(e, "=") {
				return errors.New(`execenv: each item of the first argument must be the form "key=value"`)
			}
		}
		cmd.Env = append(os.Environ(), env...)
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

func commandArgs(name string) func([]ast.Expr, *sqlx.DB) error {
	return func(args []ast.Expr, _ *sqlx.DB) error {
		list, err := toSlice(args[0].(ast.List))
		if err != nil {
			return errors.Wrap(err, name)
		}
		return stdCmd(name, list...).Run()
	}
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

func goCommand1() func([]ast.Expr, *sqlx.DB) error {
	return func(args []ast.Expr, _ *sqlx.DB) error {
		name := "go"
		cmdArgs, err := toSlice(args[1].(ast.List))
		if err != nil {
			return errors.Wrap(err, name)
		}
		var cmd *exec.Cmd
		switch lit := args[0].(*ast.Ident).Lit; lit {
		case "command":
			cmd = stdCmd(name, cmdArgs...)
		case "testall":
			// I can't be confident in using such
			// a subcommand-specific way.  Another suggestion might
			// be like `go test all`, where 'all' is a postfix
			// operator of './...'.
			cmd = stdCmd(name, append([]string{"test"}, append(cmdArgs, "./...")...)...)
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
	Fn:     goCommand1(),
}

var stackCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("stack"),
}

var leinCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("lein"),
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

var removeCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn:     remove,
}

func remove(exprs []ast.Expr, _ *sqlx.DB) error {
	s := exprs[0].(*ast.String).Lit
	fmt.Printf("remove %s?\n", s)
	for {
		fmt.Println("type y to continue")
		var ans string
		fmt.Scanf("%s", &ans)
		switch ans {
		case "i":
			fi, err := os.Stat(s)
			if err != nil {
				return err
			}
			fmt.Printf("size: %d bytes\n", fi.Size())
			fmt.Println("is directory?:", fi.IsDir())
		case "s":
			f, err := os.Open(s)
			if err != nil {
				return err
			}
			_, err = io.Copy(os.Stdout, f)
			f.Close()
			if err != nil {
				return err
			}
		case "y":
			return os.Remove(s)
		default:
			return nil
		}
	}
}

var cnpCommand = typed.Command{
	Params: []types.Type{types.String},
	Fn: func(e []ast.Expr, _ *sqlx.DB) error {
		lit := e[0].(*ast.String).Lit
		return stdCmd("create-new-project", lit).Run()
	},
}

var gvmnCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("gvmn"),
}

var vvmnCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("vvmn"),
}

var makeCommand = typed.Command{
	Params: []types.Type{types.StringList},
	Fn:     commandArgs("make"),
}

var ocamlCommand = typed.Command{
	Params: []types.Type{types.StringList},
	Fn:     commandArgs("ocaml"),
}

var rustupCommand = typed.Command{
	Params: []types.Type{types.Ident, types.StringList},
	Fn:     commandsInCommand("rustup"),
}
