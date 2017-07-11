package complete

import (
	"reflect"
	"testing"
)

func TestFile(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{
			input: "testdata/file",
			want:  []string{"1.txt", "2.txt"},
		},
		{
			input: "testdata/",
			want:  []string{"dir1/", "file1.txt", "file2.txt"},
		},
	}
	for _, test := range tests {
		list, err := File([]rune(test.input), len(test.input))
		if err != nil {
			t.Errorf("File: %v", err)
		}
		if !reflect.DeepEqual(list, test.want) {
			t.Errorf("File: want %v, got %v", test.want, list)
		}
	}
}
