package eval

import (
	"context"
	"io/ioutil"
	"testing"
)

func BenchmarkEcho(b *testing.B) {
	for i := 0; i < b.N; i++ {
		echo(context.TODO(), stream{out: ioutil.Discard}, nil, []string{"aaaaaaa"})
	}
}
