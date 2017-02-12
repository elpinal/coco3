package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/elpinal/coco3/eval"
	"github.com/elpinal/coco3/parser"
	"github.com/mattn/go-runewidth"
)

var prompt = "; "

type history struct {
	lines [][]rune
	i     int
}

func main() {
	var command = flag.String("c", "", "take first argument as a command to execute")
	flag.Parse()
	if *command != "" {
		if err := execute([]byte(*command)); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}
	var hist history
	for {
		if err := loop(&hist); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

const (
	CharCtrlA = iota + 1
	CharCtrlB
	CharCtrlC
	CharCtrlD
	CharCtrlE
	CharCtrlF
	CharCtrlG
	CharCtrlH
	CharCtrlI
	CharCtrlJ
	CharCtrlK
	CharCtrlL
	CharCtrlM
	CharCtrlN
	CharCtrlO
	CharCtrlP
	CharCtrlQ
	CharCtrlR
	CharCtrlS
	CharCtrlT
	CharCtrlU
	CharCtrlV
	CharCtrlW
	CharCtrlX
	CharCtrlY
	CharCtrlZ
	CharEscape

	CharBackspace = 127
)

func loop(hist *history) error {
	fmt.Print(prompt)
	rd := NewReader(os.Stdin)
	lines := make([][]rune, len(hist.lines))
	copy(lines, hist.lines)
	cl := &commandline{
		w: bufio.NewWriter(os.Stdout),
		hist: history{
			lines: lines,
			i:     hist.i,
		},
	}
	state, err := enterRowMode()
	if err != nil {
		return err
	}
LOOP:
	for {
		ch, err := rd.Read()
		if err != nil {
			return err
		}
		switch ch {
		case CharCtrlM:
			break LOOP
		case CharCtrlB:
			cl.moveBackward()
		case CharCtrlF:
			cl.moveForward()
		case CharBackspace:
			cl.deleteChar()
		case CharCtrlP:
			cl.prevHistory()
		case CharCtrlN:
			cl.nextHistory()
		default:
			cl.appendChar(ch)
		}
		cl.refresh()
	}
	cl.w.WriteByte('\n')
	cl.w.Flush()
	err = exitRowMode(state)
	if err != nil {
		return err
	}
	hist.lines = append(hist.lines, cl.buf)
	hist.i++
	return execute([]byte(string(cl.buf)))
}

func execute(b []byte) error {
	f, err := parser.ParseSrc(b)
	if err != nil {
		return err
	}
	return eval.Eval(f.Lines)
}

type commandline struct {
	w     *bufio.Writer
	buf   []rune
	index int

	hist history
}

func (cl *commandline) refresh() {
	cl.w.WriteString("\r\033[J")
	cl.w.WriteString(prompt)
	cl.w.WriteString(string(cl.buf))
	cl.w.WriteString(fmt.Sprintf("\033[%vG", len(prompt)+runewidth.StringWidth(string(cl.buf[:cl.index]))+1))
	cl.w.Flush()
}

func (cl *commandline) moveBackward() {
	if cl.index == 0 {
		return
	}
	cl.index--
}

func (cl *commandline) moveForward() {
	if cl.index == len(cl.buf) {
		return
	}
	cl.index++
}

func (cl *commandline) appendChar(ch rune) {
	switch cl.index {
	case 0:
		cl.buf = append([]rune{ch}, cl.buf...)
	case len(cl.buf):
		cl.buf = append(cl.buf, ch)
	default:
		buf := make([]rune, cl.index, len(cl.buf)+1)
		copy(buf, cl.buf[:cl.index])
		buf = append(buf, ch)
		buf = append(buf, cl.buf[cl.index:]...)
		cl.buf = buf
	}
	cl.index++
}

func (cl *commandline) deleteChar() {
	switch cl.index {
	case 0:
		return
	case 1:
		cl.buf = cl.buf[1:]
	case len(cl.buf):
		cl.buf = cl.buf[:cl.index-1]
	default:
		cl.buf = append(cl.buf[:cl.index-1], cl.buf[cl.index:]...)
	}
	cl.index--
}

func (cl *commandline) prevHistory() {
	if cl.hist.i == 0 {
		return
	}
	if cl.hist.i == len(cl.hist.lines) {
		cl.hist.lines = append(cl.hist.lines, cl.buf)
	} else {
		cl.hist.lines[cl.hist.i] = cl.buf
	}
	cl.hist.i--
	cl.buf = cl.hist.lines[cl.hist.i]
	cl.index = len(cl.buf)
}

func (cl *commandline) nextHistory() {
	if cl.hist.i >= len(cl.hist.lines)-1 {
		return
	}
	cl.hist.lines[cl.hist.i] = cl.buf
	cl.hist.i++
	cl.buf = cl.hist.lines[cl.hist.i]
	cl.index = len(cl.buf)
}
