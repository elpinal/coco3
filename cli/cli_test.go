package cli

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/editor"
)

func TestFlagC(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"-c", "echo aaa"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestArgs(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"testdata/basic.coco"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "aaa\nbbb\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestExit(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"-c", "echo aaa; exit 42"}
	code := c.Run(args)
	if code != 42 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	if got, want := out.String(), "aaa\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestStartUpCommand(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
		Config: &config.Config{
			StartUpCommand: []byte("echo startup..."),
		},
	}
	args := []string{"-c", "echo aaa; echo bbb"}
	code := c.Run(args)
	if code != 0 {
		t.Errorf("Run: got %v, want %v", code, 0)
	}
	if got, want := out.String(), "startup...\naaa\nbbb\n"; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestExitInStartUp(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
		Config: &config.Config{
			StartUpCommand: []byte("exit 21"),
		},
	}
	args := []string{"-c", "echo aaa; exit 42"}
	code := c.Run(args)
	if code != 21 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	if got, want := out.String(), ""; got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func TestExitInFiles(t *testing.T) {
	var out, err bytes.Buffer
	c := CLI{
		Out: &out,
		Err: &err,
	}
	args := []string{"testdata/exit1.coco", "testdata/exit2.coco"}
	code := c.Run(args)
	if code != 42 {
		t.Errorf("Run: got %v, want %v", code, 42)
	}
	golden, e := ioutil.ReadFile("testdata/exit1.golden")
	if e != nil {
		t.Errorf("reading a golden file: %v", e)
	}
	if got, want := out.String(), string(golden); got != want {
		t.Errorf("output: got %q, want %q", got, want)
	}
	if e := err.String(); e != "" {
		t.Errorf("error: %v", e)
	}
}

func BenchmarkCLI(b *testing.B) {
	c := CLI{Out: ioutil.Discard}
	for i := 0; i < b.N; i++ {
		_ = c.Run([]string{"-c", "echo 1024; echo 2048"})
	}
}

func TestCompareRunes(t *testing.T) {
	tests := []struct {
		r    []rune
		s    []rune
		want bool
	}{
		{
			r:    []rune(""),
			s:    []rune(""),
			want: true,
		},
		{
			r:    []rune("a"),
			s:    []rune("a"),
			want: true,
		},
		{
			r:    []rune("abc"),
			s:    []rune("abc"),
			want: true,
		},
		{
			r:    []rune(""),
			s:    []rune("a"),
			want: false,
		},
		{
			r:    []rune("a"),
			s:    []rune("ab"),
			want: false,
		},
		{
			r:    []rune("abc"),
			s:    []rune("aba"),
			want: false,
		},
		{
			r:    []rune("aaaaa"),
			s:    []rune("aaaab"),
			want: false,
		},
	}
	for i, test := range tests {
		got := compareRunes(test.r, test.s)
		if got != test.want {
			t.Errorf("compareRunes/%d: want %v, got %v", i, test.want, got)
		}
	}
}

func TestSanitizeHistory(t *testing.T) {
	history := []string{
		"",
		"a",
		"a",
		"",
		"a",
		"b",
		"c",
		"c",
		"",
		"b",
	}
	histRunes := sanitizeHistory(history)
	want := [][]rune{
		[]rune("a"),
		[]rune("b"),
		[]rune("c"),
		[]rune("b"),
	}
	if !reflect.DeepEqual(histRunes, want) {
		t.Errorf("want %v, got %v", want, histRunes)
	}
}

func BenchmarkRun(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := CLI{
			Config: &config.Config{
				Extra: true,
			},
		}
		_ = c.Run([]string{"-c", " "})
	}
}

type benchmarkReader struct {
	ch chan struct{}
	r  *strings.Reader
}

func (r *benchmarkReader) Read(p []byte) (n int, err error) {
	r.ch <- struct{}{}
	return r.r.Read(p)
}

func BenchmarkStartup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		ch := make(chan struct{})
		c := CLI{
			In: &benchmarkReader{
				ch: ch,
				r:  strings.NewReader(string(editor.CharEscape) + ":q" + string(editor.CharCtrlM)),
			},
			Config: &config.Config{
				HistFile: ":memory:",
			},
		}
		done := make(chan struct{})

		go func() {
			b.StartTimer()
			c.Run(nil)
			done <- struct{}{}
		}()
		<-ch
		b.StopTimer()
		<-done
		b.StartTimer()
	}
}

func TestStartup(t *testing.T) {
	c := CLI{
		In: strings.NewReader(string(editor.CharEscape) + ":q" + string(editor.CharCtrlM)),
		Config: &config.Config{
			HistFile: ":memory:",
		},
	}
	n := c.Run(nil)
	if want := 0; n != want {
		t.Errorf("Run = %d, want %d", n, want)
	}
}
