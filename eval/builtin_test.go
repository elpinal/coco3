package eval

import (
	"context"
	"io/ioutil"
	"testing"
)

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo(context.TODO(), info{
			stream: stream{out: ioutil.Discard},
			args:   []string{"aaaaaaa"},
		})
	}
}
