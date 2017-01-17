package main

import (
	"bufio"
	"bytes"
	"testing"
)

func BenchmarkCLRefresh(b *testing.B) {
	var buf bytes.Buffer
	cl := commandline{w: bufio.NewWriter(&buf)}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cl.refresh()
	}
}
