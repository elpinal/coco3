package main

import (
	"os"
	"strings"
)

func init() {
	home := os.Getenv("HOME")
	paths := []string{
		home + "/bin",
		home + "/.vvmn/vim/current/bin",
		"/usr/local/bin",
		"/usr/local/opt/coreutils/libexec/gnubin",
	}
	setPath(paths...)
}

func setPath(args ...string) {
	s := os.Getenv("PATH")
	paths := strings.Split(s, ":")
	var newPaths []string
	for _, arg := range args {
		for _, path := range paths {
			if path != arg {
				newPaths = append(newPaths, path)
			}
		}
	}
	newPaths = append(args, newPaths...)
	os.Setenv("PATH", strings.Join(newPaths, ":"))
}
