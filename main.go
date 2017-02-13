package main

// TODO: refactor

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
			case 'a':
				cl.moveForward()
				cl.mode = insertMode
			case 'i':
				cl.mode = insertMode
			case 'x':
				cl.deleteUnder()
			case 'c', 'd', 'y':
				cl.mode = operatorPendingMode
				cl.pending = ch
			case 'g':
				ch1, err := rd.Read()
				if err != nil {
					return err
				}
				switch ch1 {
				case '_':
					cl.toTheLastNonBlank()
				}
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
			case CharCtrlI:
				cl.complete()
				err := cl.completeMode(rd)
				if err != nil {
					return err
				}
			default:
				cl.appendChar(ch)
			}
		case operatorPendingMode:
			switch ch {
			case 'h':
				from := cl.index
				cl.toLeft()
				to := cl.index
				cl.operate(cl.pending, from, to)
			case 'l':
				from := cl.index
				cl.toRight()
				to := cl.index
				cl.operate(cl.pending, from, to)
			case cl.pending:
				from := 0
				to := len(cl.buf)
				cl.operate(cl.pending, from, to)
			default:
				cl.mode = normalMode
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

	mode    mode
	pending rune

	compl     []string
	maxLength int
	selected  int
	trigger   []rune
}

func (cl *commandline) construct(commands []string, width int) (omit, line int) {
	maxCount := (width + 1) / (cl.maxLength + 1)
	b := make([]byte, 0, 16)
	for i, c := range commands {
		if len(c) > width {
			omit++
			continue
		}
		if i == cl.selected-1 {
			b = append(b, "\033[7m"...)
		}
		b = append(b, c...)
		if i == cl.selected-1 {
			b = append(b, "\033[0m"...)
		}
		b = append(b, strings.Repeat(" ", cl.maxLength+1-len(c))...)
		if (i+1)%maxCount == 0 {
			line++
			b = append(b, '\n')
			cl.w.Write(b)
			b = b[:0]
		}
	}
	cl.w.Write(b)
	return omit, line
}

func (cl *commandline) refresh() {
	if cl.compl != nil && len(cl.compl) < 100 {
		width, _ := getWidth()
		cl.w.WriteString("\n\n")
		cl.w.WriteString("\r\033[J")
		_, line := cl.construct(cl.compl, int(width))
		fmt.Fprintf(cl.w, "\033[%vA", 2+line)
	} else {
		cl.w.WriteString("\r\033[J")
	}

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

func (cl *commandline) delete(from, to int) {
	left := from
	right := to
	if from > to {
		left = to
		right = from
	}
	switch {
	case left == 0:
		cl.buf = cl.buf[right:]
	case right == len(cl.buf):
		cl.buf = cl.buf[left:]
	default:
		cl.buf = append(cl.buf[:left], cl.buf[right:]...)
	}
	switch {
	case cl.index < left:
	case right < cl.index:
		cl.index = cl.index - (right - left)
	default:
		cl.index = left
	}
}

func (cl *commandline) operate(pending rune, from, to int) {
	switch pending {
	case 'c', 'd':
		cl.delete(from, to)
		if pending == 'c' {
			cl.mode = insertMode
			return
		}
	case 'y':
	}
	cl.mode = normalMode
}

func (cl *commandline) complete() {
	cl.selected = 0
	var maxLength int
	paths := strings.Split(os.Getenv("PATH"), ":")
	var compl []string
	for _, path := range paths {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			return // err
		}
		for _, file := range files {
			if strings.HasPrefix(file.Name(), string(cl.buf[:cl.index])) {
				compl = append(compl, file.Name())
				if len(file.Name()) > maxLength {
					maxLength = len(file.Name())
				}
			}
		}
	}
	cl.compl = compl
	cl.maxLength = maxLength
	cl.trigger = cl.buf[:cl.index]
}

func (cl *commandline) insert(s []rune, at int) {
	switch at {
	case 0:
		cl.buf = append(s, cl.buf...)
	case len(cl.buf):
		cl.buf = append(cl.buf, s...)
	default:
		s = append(cl.buf[:at], s...)
		cl.buf = append(s, cl.buf[at:]...)
	}
	if at <= cl.index {
		cl.index += len(s)
	}
}

func (cl *commandline) move(to int) {
	switch {
	case to >= len(cl.buf):
		switch cl.mode {
		case insertMode:
			cl.index = len(cl.buf)
		default:
			cl.index = len(cl.buf) - 1
		}
	case to <= 0:
		cl.index = 0
	default:
		cl.index = to
	}
}

func (cl *commandline) completeMode(rd *Reader) error {
	defer func() { cl.compl = nil }()
	for {
		cl.refresh()
		ch, err := rd.Read()
		if err != nil {
			return err
		}
		switch ch {
		case CharCtrlM:
			return nil
		case CharCtrlN, CharCtrlI:
			if cl.selected != len(cl.compl) {
				cl.selected++
				cl.buf = []rune(cl.compl[cl.selected-1])
				cl.index = len(cl.buf)
			} else {
				cl.selected = 0
				cl.buf = cl.trigger
				cl.index = len(cl.buf)
			}
		case CharCtrlP:
			if cl.selected != 1 {
				cl.selected--
				cl.buf = []rune(cl.compl[cl.selected-1])
				cl.index = len(cl.buf)
			} else {
				cl.selected = len(cl.compl) + 1
				cl.buf = cl.trigger
				cl.index = len(cl.buf)
			}
		case ' ':
			cl.insert([]rune{' '}, cl.index)
			return nil
		case CharBackspace:
			cl.delete(cl.index, cl.index-1)
			return nil
		default:
			cl.insert([]rune{ch}, cl.index)
			return nil
		}
	}
}
