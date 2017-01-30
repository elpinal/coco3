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
	for _, path := range paths {
		if contains(args, path) {
			continue
		}
		newPaths = append(newPaths, path)
	}
	newPaths = append(args, newPaths...)
	os.Setenv("PATH", strings.Join(newPaths, ":"))
}

func contains(x []string, s string) bool {
	for i := range x {
		if x[i] == s {
			return true
		}
	}
	return false
}
