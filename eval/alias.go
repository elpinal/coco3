package eval

import "strings"

type alias struct {
	cmd  string
	args []string
}

var aliases = make(map[string]alias)

func DefAlias(name, arg string) {
	// TODO: Support more complex syntax.
	a := strings.Split(arg, " ")
	cmd := a[0]
	args := a[1:]
	for i := range args {
		if args[i] == "''" {
			args[i] = ""
		}
	}
	if x, ok := aliases[cmd]; ok {
		cmd = x.cmd
		args = append(x.args, args...)
	}
	aliases[name] = alias{cmd, args}
}
