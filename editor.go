package main

type editor struct {
	basicEditor
	registers
}

func (e *editor) yank(r rune, from, to int) {
	s := e.slice(from, to)
	e.register(r, s)
}

func (e *editor) put(r rune, at int) {
	s := e.read(r)
	e.insert(s, at)
}
