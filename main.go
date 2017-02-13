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
	fmt.Print("\r\033[J")
	fmt.Print("\n")
	fmt.Print("-- INSERT --")
	fmt.Print("\033[A")
	fmt.Print("\r")
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
		mode: insertMode,
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
		switch cl.mode {
		case normalMode:
			switch ch {
			case CharCtrlM:
				break LOOP
			case 'h':
				cl.toLeft()
			case 'l':
				cl.toRight()
			case '0':
				cl.toTheFirst()
			case '^':
				cl.toTheFirstNonBlank()
			case '$':
				cl.toTheEnd()
			case 'a':
				cl.moveForward()
				cl.mode = insertMode
			case 'i':
				cl.mode = insertMode
			case 'x':
				cl.deleteUnder()
			case 'k':
				cl.prevHistory()
			case 'j':
				cl.nextHistory()
			case 'w':
				cl.wordForward()
			case 'W':
				cl.wordForwardNonBlank()
			case 'b':
				cl.wordBackward()
			case 'B':
				cl.wordBackwardNonBlank()
			default:
			}
		case insertMode:
			switch ch {
			case CharCtrlM:
				break LOOP
			case CharEscape:
				cl.moveBackward()
				cl.mode = normalMode
			case CharBackspace:
				cl.deleteChar()
			default:
				cl.appendChar(ch)
			}
		}
		cl.refresh()
	}
	cl.w.WriteByte('\n')
	cl.w.WriteString("\r\033[K")
	cl.w.Flush()
	err = exitRowMode(state)
	if err != nil {
		return err
	}
	if len(hist.lines) == 0 || string(hist.lines[hist.i-1]) != string(cl.buf) {
		hist.lines = append(hist.lines, cl.buf)
		hist.i++
	}
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

	mode mode
}

func (cl *commandline) refresh() {
	if cl.mode == insertMode {
		cl.w.WriteString("\033[B")
		cl.w.WriteString("\r\033[K")
		cl.w.WriteString("-- INSERT --")
		cl.w.WriteString("\033[A")
		cl.w.WriteString("\r\033[K")
	} else {
		cl.w.WriteString("\r\033[J")
	}

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
