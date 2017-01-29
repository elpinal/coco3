package eval

import "os"

func echo(args []string) error {
	if len(args) == 0 {
		_, err := os.Stdout.Write([]byte{'\n'})
		return err
	}
	_, err := os.Stdout.WriteString(args[0])
	if err != nil {
		return err
	}
	if len(args) == 1 {
		_, err := os.Stdout.Write([]byte{'\n'})
		return err
	}
	args = args[1:]
	for _, arg := range args {
		_, err := os.Stdout.Write([]byte{' '})
		if err != nil {
			return err
		}
		_, err = os.Stdout.WriteString(arg)
		if err != nil {
			return err
		}
	}
	_, err = os.Stdout.Write([]byte{'\n'})
	return err
}
