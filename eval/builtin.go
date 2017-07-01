package eval

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
)

type info struct {
	stream
	env    []string
	exitCh chan int
	args   []string
	db     *sqlx.DB
}

type stream struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

var builtins map[string]func(context.Context, info) error

func init() {
	builtins = map[string]func(context.Context, info) error{
		"cd":      cd,
		"echo":    echo,
		"exit":    exit,
		"setenv":  setenv,
		"setpath": setpath,
		"let":     let,
		"exec":    execCmd,
		"history": history,
	}
}

func cd(_ context.Context, ci info) error {
	var dir string
	switch len(ci.args) {
	case 0:
		dir = os.Getenv("HOME")
	case 1:
		dir = ci.args[0]
	default:
		return errors.New("too many arguments")
	}
	return os.Chdir(dir)
}

func echo(ctx context.Context, ci info) error {
	if len(ci.args) == 0 {
		_, err := ci.out.Write([]byte{'\n'})
		return err
	}
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	_, err := io.WriteString(ci.out, ci.args[0])
	if err != nil {
		return err
	}
	if len(ci.args) == 1 {
		_, err := ci.out.Write([]byte{'\n'})
		return err
	}
	args := ci.args[1:]
	for _, arg := range args {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		_, err := ci.out.Write([]byte{' '})
		if err != nil {
			return err
		}
		_, err = io.WriteString(ci.out, arg)
		if err != nil {
			return err
		}
	}
	_, err = ci.out.Write([]byte{'\n'})
	return err
}

func exit(_ context.Context, ci info) error {
	var code int
	switch len(ci.args) {
	case 0:
		//FIXME: exit with no args should use the exit code of the last command executed.
		code = 0
	case 1:
		i, err := strconv.Atoi(ci.args[0])
		if err != nil {
			return err
		}
		code = i
	default:
		return errors.New("too many arguments")
	}
	ci.exitCh <- code
	return nil
}

func setenv(_ context.Context, ci info) error {
	if len(ci.args)%2 == 1 {
		return errors.New("need even arguments")
	}
	for i := 0; i < len(ci.args); i += 2 {
		os.Setenv(ci.args[i], ci.args[i+1])
	}
	return nil
}

func setpath(_ context.Context, ci info) error {
	switch len(ci.args) {
	case 0:
		return errors.New("need 1 or more arguments")
	}
	paths := strings.Split(os.Getenv("PATH"), ":")
	var newPaths []string
	for _, path := range paths {
		if contains(ci.args, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	newPaths = append(ci.args, newPaths...)
	os.Setenv("PATH", strings.Join(newPaths, ":"))
	return nil
}

func let(ctx context.Context, ci info) error {
	n := getIndex(ci.args, "in")
	if n < 0 {
		return errors.New("expecting 'in', but not found")
	}
	if n == len(ci.args)-1 {
		return errors.New("expecting command name after 'in'")
	}
	if n%2 == 1 {
		return errors.New("'let ... in' should have even number of arguments")
	}
	newEnv := make([]string, 0, n/2)
	for i := 0; i < n; i += 2 {
		newEnv = append(newEnv, ci.args[i]+"="+ci.args[i+1])
	}
	name := ci.args[n+1]
	// Builtin command is not supported.
	// For instance, `let ... in cd` does actually execute /usr/bin/cd.
	cmd := exec.CommandContext(ctx, name, ci.args[n+2:]...)
	cmd.Env = append(ci.env, newEnv...)
	cmd.Stdin = ci.in
	cmd.Stdout = ci.out
	cmd.Stderr = ci.err
	return cmd.Run()
}

func getIndex(x []string, s string) int {
	for i := range x {
		if x[i] == s {
			return i
		}
	}
	return -1
}

func contains(x []string, s string) bool {
	for i := range x {
		if x[i] == s {
			return true
		}
	}
	return false
}

func execCmd(ctx context.Context, ci info) error {
	if len(ci.args) == 0 {
		return errors.New("1 or more arguments required")
	}
	name, err := exec.LookPath(ci.args[0])
	if err != nil {
		return err
	}
	return syscall.Exec(name, append([]string{name}, ci.args[1:]...), ci.env)
}

type execution struct {
	Time time.Time
	Line string
}

func history(ctx context.Context, ci info) error {
	buf := bufio.NewWriter(ci.out)
	data := execution{}
	rows, err := ci.db.Queryx("select * from command_info")
	if err != nil {
		return err
	}
	for rows.Next() {
		err := rows.StructScan(&data)
		if err != nil {
			return err
		}
		fmt.Fprintf(buf, "%v  %s\n", data.Time.Format("Mon, 02 Jan 2006 15:04:05"), data.Line)
	}
	return buf.Flush()
}
