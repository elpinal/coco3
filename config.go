package main

import "os"

// This file is *MY* config file.
// TODO: Remove personal setting later.
func init() {
	home := os.Getenv("HOME")
	os.Setenv("GHQ_ROOT", home+"/src")
	os.Setenv("GOPATH", home)
	os.Setenv("EDITOR", "vim")
	os.Setenv("PAGER", "less")
}
