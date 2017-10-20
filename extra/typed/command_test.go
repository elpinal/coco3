package typed

import (
	"testing"

	"github.com/elpinal/coco3/extra/types"
)

func TestSignature(t *testing.T) {
	cmd := Command{Params: []types.Type{types.Int, types.String}}
	got := string(cmd.Signature())
	want := "Int -> String"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
