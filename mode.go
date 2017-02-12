package main

type mode int

const (
	normalMode mode = iota + 1
	visualMode
	selectMode
	insertMode
	commandlineMode
	exMode
)
