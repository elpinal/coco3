package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var prompt = "Σ> "

func main() {
	var err error
	for err == nil {
		err = loop()
	}
	fmt.Fprintln(os.Stderr, err)
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

func loop() error {
	fmt.Print(prompt)
	rd := NewReader(os.Stdin)
	cl := &commandline{
		w: bufio.NewWriter(os.Stdout),
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
		default:
			cl.buf = append(cl.buf, ch)
			cl.index++
		}
		cl.refresh()
	}
	cl.w.WriteByte('\n')
	cl.w.Flush()
	err = exitRowMode(state)
	if err != nil {
		return err
	}
	cmd := exec.Command(string(cl.buf))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

type commandline struct {
	w     *bufio.Writer
	buf   []rune
	index int
}

func (cl *commandline) refresh() {
	cl.w.WriteString("\r\033[J")
	cl.w.WriteString(prompt)
	cl.w.WriteString(string(cl.buf))
	cl.w.WriteString(fmt.Sprintf("\033[%vG", len(prompt)+cl.index+1))
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
