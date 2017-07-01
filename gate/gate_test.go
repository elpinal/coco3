package gate

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/elpinal/coco3/config"
	"github.com/elpinal/coco3/editor"
)

func TestGate(t *testing.T) {
	in := strings.NewReader("echo 1" + string(editor.CharCtrlM) + string(editor.CharEscape) + "ka" + string(editor.CharBackspace) + "2" + string(editor.CharCtrlM))
	conf := new(config.Config)
	conf.Init()
	g := New(conf, in, ioutil.Discard, ioutil.Discard, nil).(*gate)
	for _, n := range [...]string{"1", "2"} {
		b, err := g.Read()
		if err != nil {
			t.Errorf("reading input: %v", err)
		}
		if want := "echo " + n; string(b) != want {
			t.Errorf("got %q, want %q", string(b), want)
		}
	}
	if l := len(g.history); l != 2 {
		t.Errorf("the length of history should be %v, got %v", 2, l)
	}
}
