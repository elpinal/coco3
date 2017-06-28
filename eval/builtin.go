package eval

import (
	"context"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

type stream struct {
	in  io.Reader
	out io.Writer
	err io.Writer
}

var builtins map[string]func(context.Context, stream, []string, *Evaluator, []string) error

func init() {
	builtins = map[string]func(context.Context, stream, []string, *Evaluator, []string) error{
		"cd":      cd,
		"echo":    echo,
		"exit":    exit,
		"setenv":  setenv,
		"setpath": setpath,
		"let":     let,
	}
}

func cd(_ context.Context, _ stream, env []string, _ *Evaluator, args []string) error {
	var dir string
	switch len(args) {
	case 0:
		dir = os.Getenv("HOME")
	case 1:
		dir = args[0]
	default:
		return errors.New("too many arguments")
	}
	return os.Chdir(dir)
}

func echo(ctx context.Context, s stream, env []string, _ *Evaluator, args []string) error {
	if len(args) == 0 {
		_, err := s.out.Write([]byte{'\n'})
		return err
	}
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	_, err := io.WriteString(s.out, args[0])
	if err != nil {
		return err
	}
	if len(args) == 1 {
		_, err := s.out.Write([]byte{'\n'})
		return err
	}
	args = args[1:]
	for _, arg := range args {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		_, err := s.out.Write([]byte{' '})
		if err != nil {
			return err
		}
		_, err = io.WriteString(s.out, arg)
		if err != nil {
			return err
		}
	}
	_, err = s.out.Write([]byte{'\n'})
	return err
}

func exit(_ context.Context, _ stream, env []string, e *Evaluator, args []string) error {
	var code int
	switch len(args) {
	case 0:
		//FIXME: exit with no args should use the exit code of the last command executed.
		code = 0
	case 1:
		i, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		code = i
	default:
		return errors.New("too many arguments")
	}
	e.ExitCh <- code
	return nil
}

func setenv(_ context.Context, _ stream, env []string, _ *Evaluator, args []string) error {
	if len(args)%2 == 1 {
		return errors.New("need even arguments")
	}
	for i := 0; i < len(args); i += 2 {
		os.Setenv(args[i], args[i+1])
	}
	return nil
}

func setpath(_ context.Context, _ stream, env []string, _ *Evaluator, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("need 1 or more arguments")
	}
	paths := strings.Split(os.Getenv("PATH"), ":")
	var newPaths []string
	for _, path := range paths {
		if contains(args, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	newPaths = append(args, newPaths...)
	os.Setenv("PATH", strings.Join(newPaths, ":"))
	return nil
}

func let(ctx context.Context, s stream, env []string, e *Evaluator, args []string) error {
	n := getIndex(args, "in")
	if n < 0 {
		return errors.New("expecting 'in', but not found")
	}
	if n == len(args)-1 {
		return errors.New("expecting command name after 'in'")
	}
	if n%2 == 1 {
		return errors.New("'let ... in' should have even number of arguments")
	}
	newEnv := make([]string, 0, n/2)
	for i := 0; i < n; i += 2 {
		newEnv = append(newEnv, args[i]+"="+args[i+1])
	}
	name := args[n+1]
	cmd := e.CommandContext(ctx, name, args[n+2:]...)
	cmd.SetEnv(append(newEnv, env...))
	cmd.SetStdin(s.in)
	cmd.SetStdout(s.out)
	cmd.SetStderr(s.err)
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
