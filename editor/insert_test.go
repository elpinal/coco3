package editor

import (
	"fmt"
	"strings"
	"testing"
)

func testInsertMode(input string, tt []string) func(*testing.T) {
	return func(t *testing.T) {
		i := insert{}
		i.init()
		i.in = NewReader(strings.NewReader(input))
		for n := range input {
			c, _, err := i.Run()
			if err != nil {
				t.Errorf("Run: %v", err)
			}
			if c != cont {
				t.Errorf("Run: want %v, but got %v", cont, c)
			}
			if got, want := string(i.Runes()), tt[n]; got != want {
				t.Errorf("Run/%d: want %q, but got %q", n, want, got)
			}
		}
	}
}

func TestInputMatches(t *testing.T) {
	tests := []struct {
		input  string
		expect []string
	}{
		{
			input:  "abc",
			expect: []string{"a", "ab", "abc"},
		},
		{
			input:  "a 'b",
			expect: []string{"a", "a ", "a ''", "a 'b'"},
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("%d", n), testInsertMode(test.input, test.expect))
	}
}
