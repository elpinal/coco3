package editor

import (
	"strings"
	"testing"
)

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		command := strings.Repeat("hjkl", 100)
		ed := newNormal(
			streamSet{
				in: NewReader(strings.NewReader(command)),
			},
			newEditor(),
		)
		for range command {
			_, _, err := ed.Run()
			if err != nil {
				b.Errorf("normal: %v", err)
			}
		}
	}
}

func TestUpdateNumber(t *testing.T) {
	norm := newNormal(
		streamSet{
			in: NewReader(nil),
		},
		newEditor(),
	)
	norm.updateNumber(func(n int) int {
		return n
	})
}
