package eval

import "strings"

func init() {
	// TODO: Remove this.
	DefAlias("ls", "ls --show-control-chars --color=auto")
	DefAlias("la", "ls -a")
	DefAlias("ll", "ls -l")
	DefAlias("lla", "ls -la")
	DefAlias("v", "vim")
	DefAlias("g", "git")
}

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
	if x, ok := aliases[cmd]; ok {
		cmd = x.cmd
		args = append(x.args, args...)
	}
	aliases[name] = alias{cmd, args}
}
