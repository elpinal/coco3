package eval

import (
	"errors"
	"io"
	"os"
	"strconv"
)

var builtins = map[string]func(*Evaluator, []string) error{
	"cd":   cd,
	"echo": echo,
	"exit": exit,
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
