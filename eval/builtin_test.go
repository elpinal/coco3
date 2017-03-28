package eval

import (
	"io/ioutil"
	"testing"
)

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo(stream{out: ioutil.Discard}, nil, []string{"aaaaaaa"})
	}
}
