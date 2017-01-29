package eval

import "strings"

func init() {
	defAlias("ls", "ls --show-control-chars --color=auto")
	defAlias("la", "ls -a")
	defAlias("ll", "ls -l")
	defAlias("lla", "ls -la")
	defAlias("v", "vim")
	defAlias("g", "git")
}

type alias struct {
	cmd  string
	args []string
}

var aliases = make(map[string]alias)

func defAlias(name, arg string) {
	a := strings.Split(arg, " ")
	cmd := a[0]
	args := a[1:]
	if x, ok := aliases[cmd]; ok {
		cmd = x.cmd
		args = append(x.args, args...)
	}
	aliases[name] = alias{cmd, args}
}
