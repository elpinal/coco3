package main

type mode int

const (
	normalMode mode = iota + 1
	visualMode
	selectMode
	insertMode
	commandlineMode
	exMode

	operatorPendingMode
	replaceMode
	virtualReplaceMode
	insertNormalMode
	insertVisualMode
	insertSelectMode
)
