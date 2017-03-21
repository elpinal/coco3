package eval

import (
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

var builtins = map[string]func(*Evaluator, []string) error{
	"cd":      cd,
	"echo":    echo,
	"exit":    exit,
	"setpath": setpath,
	"setenv":  setenv,
}

func cd(_ *Evaluator, args []string) error {
	var dir string
	switch len(args) {
	case 0:
		dir = os.Getenv("HOME")
	case 1:
		dir = args[0]
	default:
		return errors.New("cd: too many arguments")
	}
	return os.Chdir(dir)
}

func echo(e *Evaluator, args []string) error {
	if len(args) == 0 {
		_, err := e.out.Write([]byte{'\n'})
		return err
	}
	_, err := io.WriteString(e.out, args[0])
	if err != nil {
		return err
	}
	if len(args) == 1 {
		_, err := e.out.Write([]byte{'\n'})
		return err
	}
	args = args[1:]
	for _, arg := range args {
		_, err := e.out.Write([]byte{' '})
		if err != nil {
			return err
		}
		_, err = io.WriteString(e.out, arg)
		if err != nil {
			return err
		}
	}
	_, err = e.out.Write([]byte{'\n'})
	return err
}

func exit(_ *Evaluator, args []string) error {
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
		return errors.New("exit: too many arguments")
	}
	os.Exit(code)
	return nil
}

func setpath(_ *Evaluator, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("setpath: need 1 or more arguments")
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

func contains(x []string, s string) bool {
	for i := range x {
		if x[i] == s {
			return true
		}
	}
	return false
}

func setenv(_ *Evaluator, args []string) error {
	if len(args)%2 == 1 {
		return errors.New("setenv: need even arguments")
	}
	for i := 0; i < len(args); i += 2 {
		os.Setenv(args[i], args[i+1])
	}
	return nil
}
