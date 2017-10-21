package typed

import (
	"testing"

	"github.com/elpinal/coco3/extra/types"
)

func TestSignature(t *testing.T) {
	tests := []struct {
		ts   []types.Type
		want string
	}{
		{
			ts:   []types.Type{},
			want: "",
		},
		{
			ts:   []types.Type{types.Int},
			want: "Int",
		},
		{
			ts:   []types.Type{types.Int, types.String},
			want: "Int -> String",
		},
		{
			ts:   []types.Type{types.Int, types.String, types.StringList},
			want: "Int -> String -> List String",
		},
	}
	for i, test := range tests {
		cmd := Command{Params: test.ts}
		got := string(cmd.Signature())
		if got != test.want {
			t.Errorf("Signature/%d: got %q, want %q", i, got, test.want)
		}
	}
}
