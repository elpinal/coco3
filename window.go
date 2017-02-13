package main

import (
	"syscall"
	"unsafe"
)

type window struct {
	row    uint16
	col    uint16
}

func getWidth() (uint, error) {
	ws := &window{}
	code, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdin), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(ws)))
	if int(code) == -1 {
		return 0, errno
	}
	return uint(ws.col), nil
}
