package editor

import (
	"context"
	"strings"
	"testing"
)

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		command := strings.Repeat("hjkl", 100)
		ed := normal{
			streamSet: streamSet{in: NewReaderContext(context.TODO(), strings.NewReader(command))},
			editor:    newEditor(),
		}
		for range command {
			_, _, err := ed.Run()
			if err != nil {
				b.Errorf("normal: %v", err)
			}
		}
	}
}
